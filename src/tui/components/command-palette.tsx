import React, { useState, useMemo } from 'react';
import { Box, Text, useInput } from 'ink';
import { colors } from '../theme.js';

export interface PaletteItem {
  name: string;
  description: string;
  shortcut?: string;
}

interface CommandPaletteProps {
  items: PaletteItem[];
  onSelect: (name: string) => void;
  onCancel: () => void;
  lang?: string;
}

export function CommandPalette({
  items,
  onSelect,
  onCancel,
  lang,
}: CommandPaletteProps): React.ReactElement {
  const ko = lang === 'ko';
  const [query, setQuery] = useState('');
  const [cursor, setCursor] = useState(0);

  const filtered = useMemo(() => {
    if (query.length === 0) return items;
    const q = query.toLowerCase();
    return items.filter(
      (item) =>
        item.name.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q),
    );
  }, [items, query]);

  const safeCursor = Math.min(cursor, Math.max(0, filtered.length - 1));

  useInput((input, key) => {
    if (key.escape) {
      onCancel();
      return;
    }
    if (key.return) {
      const item = filtered[safeCursor];
      if (item) onSelect(item.name);
      return;
    }
    if (key.upArrow) {
      setCursor((prev) => (prev <= 0 ? filtered.length - 1 : prev - 1));
      return;
    }
    if (key.downArrow) {
      setCursor((prev) => (prev >= filtered.length - 1 ? 0 : prev + 1));
      return;
    }
    if (key.backspace || key.delete) {
      setQuery((prev) => prev.slice(0, -1));
      setCursor(0);
      return;
    }
    if (key.ctrl || key.meta || key.tab) return;
    if (input) {
      setQuery((prev) => prev + input);
      setCursor(0);
    }
  });

  const maxVisible = 12;
  const startIdx = Math.max(0, safeCursor - maxVisible + 1);
  const visible = filtered.slice(startIdx, startIdx + maxVisible);

  return (
    <Box
      flexDirection="column"
      borderStyle="round"
      borderColor={colors.borderFocus}
      paddingX={1}
    >
      <Box>
        <Text bold color={colors.model}>
          {ko ? '명령어 팔레트' : 'Command Palette'}
        </Text>
        <Text color={colors.dimText}> {'\u2014'} {ko ? '입력하여 검색' : 'type to filter'}</Text>
      </Box>
      <Box>
        <Text color={colors.prompt}>{'\u276F'} </Text>
        <Text color={colors.userText}>{query}</Text>
        <Text color={colors.dimText}>{'\u2588'}</Text>
      </Box>
      <Text color={colors.dimText}>{'\u2500'.repeat(40)}</Text>
      {visible.length === 0 ? (
        <Text color={colors.dimText}> {ko ? '일치하는 명령어 없음' : 'No matching commands'}</Text>
      ) : (
        visible.map((item, i) => {
          const actualIdx = startIdx + i;
          const isCursor = actualIdx === safeCursor;
          return (
            <Box key={item.name}>
              <Text color={isCursor ? colors.hint : colors.dimText}>
                {isCursor ? '\u25B8' : ' '}{' '}
              </Text>
              <Text
                color={isCursor ? colors.hint : colors.assistantText}
                bold={isCursor}
              >
                /{item.name}
              </Text>
              <Text color={colors.dimText}>
                {' \u2014 '}
                {item.description}
              </Text>
              {item.shortcut && (
                <Text color={colors.warning}> [{item.shortcut}]</Text>
              )}
            </Box>
          );
        })
      )}
      <Text color={colors.dimText}>{'\u2500'.repeat(40)}</Text>
      <Text color={colors.dimText}>
        {ko ? '\u2191\u2193 이동  Enter 실행  Esc 닫기' : '\u2191\u2193 navigate  Enter execute  Esc close'}
      </Text>
    </Box>
  );
}
