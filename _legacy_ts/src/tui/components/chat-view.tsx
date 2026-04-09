import React, { useMemo, useState, useEffect, useRef } from 'react';
import { Box, Text, useInput } from 'ink';
import { Spinner } from './spinner.js';
import { colors } from '../theme.js';
import { FileDiff, isDiffContent, extractDiffFilePath } from './file-diff.js';

export interface DisplayMessage {
  id: string;
  role: 'user' | 'assistant' | 'tool-call' | 'tool-result' | 'system';
  content: string;
  toolName?: string;
}

interface ChatViewProps {
  messages: DisplayMessage[];
  streamingText: string;
  isLoading: boolean;
  height: number;
  verbose?: boolean;
}

const COMPACT_LINES = 10;
const VERBOSE_LINES = 500;

function truncateLines(
  text: string,
  maxLines: number,
): { text: string; hidden: number } {
  const lines = text.split('\n');
  if (lines.length <= maxLines) return { text, hidden: 0 };
  const kept = lines.slice(0, maxLines);
  const hidden = lines.length - maxLines;
  return { text: kept.join('\n'), hidden };
}

/** Memoized message bubble — avoids re-rendering unchanged messages */
const MessageBubble = React.memo(function MessageBubble({
  message,
  maxToolLines,
}: {
  message: DisplayMessage;
  maxToolLines: number;
}): React.ReactElement {
  switch (message.role) {
    case 'user':
      return (
        <Box paddingX={1}>
          <Text color={colors.userText} bold>
            {'\u276F'}{' '}
          </Text>
          <Text color={colors.userText}>{message.content}</Text>
        </Box>
      );

    case 'assistant':
      return (
        <Box paddingX={1} flexDirection="column">
          <Text color={colors.assistantText}>{message.content}</Text>
        </Box>
      );

    case 'tool-call':
      return (
        <Box paddingX={1}>
          <Text color={colors.toolCall}>
            {'\u26A1'} {message.toolName ?? 'tool'}
          </Text>
        </Box>
      );

    case 'tool-result': {
      const { text: truncated, hidden } = truncateLines(
        message.content,
        maxToolLines,
      );
      const isError = message.content.startsWith('Error:');
      const showDiff = isDiffContent(message.content) &&
        (message.toolName === 'git_diff' || message.toolName === 'edit_file' ||
         message.toolName === 'patch' || message.toolName === 'hashline_edit');
      return (
        <Box paddingX={1} flexDirection="column">
          <Text color={colors.toolCall} dimColor>
            {'\u2514'} {message.toolName ?? 'tool'}
          </Text>
          <Box paddingLeft={2} flexDirection="column">
            {showDiff ? (
              <FileDiff
                diff={truncated}
                filePath={extractDiffFilePath(message.content)}
              />
            ) : (
              <Text color={isError ? colors.error : colors.toolResult}>
                {truncated}
              </Text>
            )}
            {hidden > 0 && (
              <Text color={colors.dimText}>
                {'\u2501\u2501'} {hidden} more lines (Ctrl+O verbose){' '}
                {'\u2501\u2501'}
              </Text>
            )}
          </Box>
        </Box>
      );
    }

    case 'system':
      return (
        <Box paddingX={1} flexDirection="column">
          <Text color={colors.hint}>{message.content}</Text>
        </Box>
      );

    default: {
      const _exhaustive: never = message.role;
      return <Text>{String(_exhaustive)}</Text>;
    }
  }
});

function estimateLines(
  msg: DisplayMessage,
  width: number,
  maxToolLines: number,
): number {
  const usableWidth = Math.max(width - 4, 20);
  const content =
    msg.role === 'tool-result'
      ? truncateLines(msg.content, maxToolLines).text
      : msg.content;
  let lines = 0;
  for (const line of content.split('\n')) {
    lines += Math.max(1, Math.ceil(line.length / usableWidth));
  }
  if (msg.role === 'user') lines += 1;
  if (
    msg.role === 'tool-result' &&
    msg.content.split('\n').length > maxToolLines
  ) {
    lines += 1;
  }
  return lines;
}

export const ChatView = React.memo(function ChatView({
  messages,
  streamingText,
  isLoading,
  height,
  verbose = false,
}: ChatViewProps): React.ReactElement {
  const width = process.stdout.columns || 80;
  const maxToolLines = verbose ? VERBOSE_LINES : COMPACT_LINES;

  // Scroll state: 0 = bottom (latest), positive = scrolled up by N LINES
  const [scrollOffset, setScrollOffset] = useState(0);

  // Per-message line heights (memoized)
  const lineHeights = useMemo(
    () => messages.map((msg) => estimateLines(msg, width, maxToolLines)),
    [messages, width, maxToolLines],
  );
  const totalLines = useMemo(
    () => lineHeights.reduce((a, b) => a + b, 0),
    [lineHeights],
  );

  // Auto-scroll to bottom on new messages
  const msgCount = messages.length;
  useEffect(() => { setScrollOffset(0); }, [msgCount]);

  // Auto-scroll when streaming starts
  const isStreaming = streamingText.length > 0;
  useEffect(() => { if (isStreaming) setScrollOffset(0); }, [isStreaming]);

  // Scroll with PageUp/PageDown (line-based) + Shift+Up/Down fallback
  const scrollStep = Math.max(3, Math.floor(height / 3));
  useInput((_input, key) => {
    const up = key.pageUp || (key.upArrow && (key.meta || key.shift));
    const down = key.pageDown || (key.downArrow && (key.meta || key.shift));
    if (up) {
      setScrollOffset((prev) => Math.min(prev + scrollStep, Math.max(0, totalLines - 1)));
    }
    if (down) {
      setScrollOffset((prev) => Math.max(prev - scrollStep, 0));
    }
  });

  // Mouse wheel scroll — intercept SGR mouse sequences BEFORE Ink sees them
  const totalLinesRef = useRef(totalLines);
  totalLinesRef.current = totalLines;
  const setScrollRef = useRef(setScrollOffset);
  setScrollRef.current = setScrollOffset;

  useEffect(() => {
    const stdin = process.stdin;
    type EmitFn = (event: string | symbol, ...args: unknown[]) => boolean;
    const origEmit: EmitFn = stdin.emit.bind(stdin);
    // Match mouse sequences WITH or WITHOUT leading ESC (readline may strip it)
    const mouseSeqRe = /(?:\x1b)?\[<(\d+);\d+;\d+[Mm]/g;
    const mouseStripRe = /(?:\x1b)?\[<\d+;\d+;\d+[Mm]/g;

    // Enable SGR extended mouse mode (wheel events only via button tracking)
    process.stdout.write('\x1b[?1000h\x1b[?1006h');

    // Patch stdin.emit to intercept mouse data before Ink/TextInput sees it
    const patchedEmit: EmitFn = (event, ...args) => {
      if (event === 'data' && args.length > 0) {
        const raw = args[0];
        const str = Buffer.isBuffer(raw) ? raw.toString() : typeof raw === 'string' ? raw : '';
        if (str.includes('[<')) {
          // Handle wheel events
          let match: RegExpExecArray | null;
          mouseSeqRe.lastIndex = 0;
          while ((match = mouseSeqRe.exec(str)) !== null) {
            const btn = parseInt(match[1]!, 10);
            if (btn === 64) {
              setScrollRef.current((prev) => Math.min(prev + 3, Math.max(0, totalLinesRef.current - 1)));
            } else if (btn === 65) {
              setScrollRef.current((prev) => Math.max(prev - 3, 0));
            }
          }
          // Strip ALL mouse sequences so Ink never sees them
          const cleaned = str.replace(mouseStripRe, '');
          if (!cleaned) return true; // fully consumed
          return origEmit('data', Buffer.from(cleaned));
        }
      }
      return origEmit(event, ...args);
    };

    stdin.emit = patchedEmit as NodeJS.ReadStream['emit'];

    return () => {
      stdin.emit = origEmit;
      process.stdout.write('\x1b[?1000l\x1b[?1006l');
    };
  }, []);

  // Compute visible messages with line-level precision
  const { visibleMsgs, linesAbove, linesBelow } = useMemo(() => {
    // Reserve space for streaming text, loading spinner, indicators
    let reserved = 2; // up/down scroll indicators (worst case)
    if (isStreaming) {
      for (const line of streamingText.split('\n')) {
        reserved += Math.max(1, Math.ceil(line.length / Math.max(width - 4, 20)));
      }
    }
    if (isLoading && !isStreaming) reserved += 1;

    const available = Math.max(height - reserved, 1);
    const maxScroll = Math.max(0, totalLines - available);
    const clamped = Math.min(scrollOffset, maxScroll);

    // Visible line window (absolute positions from top)
    const windowBottom = totalLines - clamped;
    const windowTop = Math.max(0, windowBottom - available);

    const result: DisplayMessage[] = [];
    let acc = 0;

    for (let i = 0; i < messages.length; i++) {
      const h = lineHeights[i] ?? 0;
      const msgTop = acc;
      const msgBottom = acc + h;
      acc = msgBottom;

      // Skip messages entirely outside the window
      if (msgBottom <= windowTop || msgTop >= windowBottom) continue;

      const msg = messages[i]!;
      const skipFromTop = Math.max(0, windowTop - msgTop);
      const skipFromBottom = Math.max(0, msgBottom - windowBottom);

      if (skipFromTop > 0 || skipFromBottom > 0) {
        // Partially visible — truncate content at line boundaries
        const rawLines = (msg.role === 'tool-result'
          ? truncateLines(msg.content, maxToolLines).text
          : msg.content
        ).split('\n');
        const start = Math.min(skipFromTop, Math.max(0, rawLines.length - 1));
        const end = Math.max(start + 1, rawLines.length - skipFromBottom);
        result.push({ ...msg, content: rawLines.slice(start, end).join('\n') });
      } else {
        result.push(msg);
      }
    }

    return { visibleMsgs: result, linesAbove: windowTop, linesBelow: clamped };
  }, [messages, streamingText, isLoading, isStreaming, height, width, scrollOffset, maxToolLines, lineHeights, totalLines]);

  return (
    <Box flexDirection="column" height={height} overflow="hidden" justifyContent="flex-end">
      {/* Scroll-up indicator */}
      {linesAbove > 0 && (
        <Box paddingX={1} justifyContent="center">
          <Text color={colors.dimText}>
            {'\u2191'} {linesAbove} lines (PgUp)
          </Text>
        </Box>
      )}

      {/* Empty state */}
      {messages.length === 0 && !isLoading && (
        <Box justifyContent="center" flexGrow={1} alignItems="center">
          <Text color={colors.dimText}>
            Type a message to start... {'\u2502'} Ctrl+C to exit
          </Text>
        </Box>
      )}

      {/* Messages — visible portion */}
      {visibleMsgs.map((msg) => (
        <MessageBubble key={msg.id} message={msg} maxToolLines={maxToolLines} />
      ))}

      {/* Active streaming text */}
      {isStreaming && (
        <Box paddingX={1} flexDirection="column" overflow="hidden">
          <Text color={colors.assistantText} wrap="wrap">{streamingText}</Text>
        </Box>
      )}

      {/* Thinking indicator */}
      {isLoading && !isStreaming && (
        <Box paddingX={1}>
          <Spinner label="Thinking..." color={colors.statusThinking} />
        </Box>
      )}

      {/* Scroll-down indicator */}
      {linesBelow > 0 && (
        <Box paddingX={1} justifyContent="center">
          <Text color={colors.dimText}>
            {'\u2193'} {linesBelow} lines (PgDn)
          </Text>
        </Box>
      )}
    </Box>
  );
})
