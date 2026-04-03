import React, { useState, useMemo } from 'react';
import { Box, Text, useInput } from 'ink';
import { colors } from '../theme.js';

export interface MenuItem {
  label: string;
  value: string;
  active?: boolean;
  badge?: string;
  badgeColor?: string;
}

interface SelectMenuProps {
  title: string;
  items: MenuItem[];
  onSelect: (value: string) => void;
  onCancel: () => void;
  legend?: string;
  /** Called when cursor moves — use for live preview (e.g. theme) */
  onHighlight?: (value: string) => void;
  /** Max visible items before scrolling (default: terminal height - 8) */
  maxVisible?: number;
}

export function SelectMenu({ title, items, onSelect, onCancel, legend, onHighlight, maxVisible }: SelectMenuProps): React.ReactElement {
  const [cursor, setCursor] = useState(0);

  // 스크롤 윈도우 크기: 터미널 높이 기반 또는 props
  const windowSize = useMemo(() => {
    const termRows = process.stdout.rows || 24;
    return maxVisible ?? Math.max(8, termRows - 8); // 헤더/푸터/테두리 공간 제외
  }, [maxVisible]);

  // 스크롤 오프셋 계산: 커서가 항상 윈도우 안에 있도록
  const scrollOffset = useMemo(() => {
    if (items.length <= windowSize) return 0;
    // 커서를 윈도우 중앙 부근에 유지
    const half = Math.floor(windowSize / 2);
    let offset = cursor - half;
    offset = Math.max(0, offset);
    offset = Math.min(items.length - windowSize, offset);
    return offset;
  }, [cursor, items.length, windowSize]);

  const visibleItems = items.slice(scrollOffset, scrollOffset + windowSize);
  const hasScrollUp = scrollOffset > 0;
  const hasScrollDown = scrollOffset + windowSize < items.length;

  useInput((input, key) => {
    if (key.escape) {
      onCancel();
      return;
    }

    if (key.upArrow) {
      setCursor((prev) => {
        const next = prev <= 0 ? items.length - 1 : prev - 1;
        if (onHighlight && items[next]) onHighlight(items[next].value);
        return next;
      });
      return;
    }

    if (key.downArrow) {
      setCursor((prev) => {
        const next = prev >= items.length - 1 ? 0 : prev + 1;
        if (onHighlight && items[next]) onHighlight(items[next].value);
        return next;
      });
      return;
    }

    if (key.return) {
      const item = items[cursor];
      if (item) onSelect(item.value);
      return;
    }

    // Number key shortcuts (0-9): 1-9 for items 1-9, 0 for item 10
    const num = parseInt(input, 10);
    if (!isNaN(num)) {
      const idx = num === 0 ? 9 : num - 1; // 0 key = 10th item
      if (idx >= 0 && idx < items.length) {
        const item = items[idx];
        if (item) onSelect(item.value);
      }
    }
  });

  return (
    <Box flexDirection="column" borderStyle="round" borderColor={colors.borderFocus} paddingX={1}>
      <Text bold color={colors.model}>{title}</Text>
      <Text color={colors.dimText}>{'\u2500'.repeat(30)}</Text>
      {hasScrollUp && (
        <Text color={colors.dimText}>  {'\u25B2'} {scrollOffset}개 더 위에</Text>
      )}
      {visibleItems.map((item, vi) => {
        const i = vi + scrollOffset; // 실제 인덱스
        const isCursor = i === cursor;
        const marker = item.active ? '\u25CF' : isCursor ? '\u25B8' : ' ';
        const shortcut = i < 26 ? String.fromCharCode(i < 9 ? 49 + i : i === 9 ? 48 : 97 + i - 10) : ' ';
        return (
          <Box key={item.value}>
            <Text color={isCursor ? colors.model : colors.dimText}>{shortcut} </Text>
            <Text color={isCursor ? colors.hint : item.active ? colors.success : colors.assistantText}>
              {marker} {item.label}
            </Text>
            {item.badge && (
              <Text color={item.badgeColor ?? colors.dimText}> {item.badge}</Text>
            )}
          </Box>
        );
      })}
      {hasScrollDown && (
        <Text color={colors.dimText}>  {'\u25BC'} {items.length - scrollOffset - windowSize}개 더 아래</Text>
      )}
      <Text color={colors.dimText}>{'\u2500'.repeat(30)}</Text>
      {items.length > windowSize && (
        <Text color={colors.dimText}>{cursor + 1}/{items.length}</Text>
      )}
      {legend && <Text color={colors.dimText}>{legend}</Text>}
      <Text color={colors.dimText}>{'↑↓ navigate  Enter select  Esc cancel'}</Text>
    </Box>
  );
}
