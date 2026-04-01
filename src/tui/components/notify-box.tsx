import React from 'react';
import { Box, Text } from 'ink';
import { colors } from '../theme.js';

type NotifyType = 'info' | 'warning' | 'success' | 'error' | 'tip';

interface NotifyBoxProps {
  type?: NotifyType;
  title?: string;
  children: React.ReactNode;
  width?: number;
}

const TYPE_STYLES: Record<NotifyType, { icon: string; color: string; border: 'single' | 'double' | 'round' }> = {
  info:    { icon: '\u2139',  color: '#89b4fa', border: 'round' },   // i
  warning: { icon: '\u26A0',  color: '#f9e2af', border: 'single' },  // warning sign
  success: { icon: '\u2714',  color: '#a6e3a1', border: 'round' },   // checkmark
  error:   { icon: '\u2718',  color: '#f38ba8', border: 'double' },  // x mark
  tip:     { icon: '\u2728',  color: '#cba6f7', border: 'round' },   // sparkles (actually \u2728 doesn't work in all terminals)
};

export const NotifyBox = React.memo(function NotifyBox({
  type = 'info',
  title,
  children,
  width,
}: NotifyBoxProps): React.ReactElement {
  const style = TYPE_STYLES[type];

  return (
    <Box
      flexDirection="column"
      borderStyle={style.border}
      borderColor={style.color}
      paddingX={1}
      width={width}
      marginX={1}
      marginY={0}
    >
      {title && (
        <Box>
          <Text color={style.color} bold>
            {style.icon} {title}
          </Text>
        </Box>
      )}
      <Box flexDirection="column">
        {typeof children === 'string' ? (
          <Text color={colors.dimText}>{children}</Text>
        ) : (
          children
        )}
      </Box>
    </Box>
  );
});

/**
 * Inline notification — single line, no border. Good for compact tips.
 */
export function NotifyInline({
  type = 'info',
  children,
}: {
  type?: NotifyType;
  children: React.ReactNode;
}): React.ReactElement {
  const style = TYPE_STYLES[type];
  return (
    <Box paddingX={2}>
      <Text color={style.color}>
        {style.icon}{' '}
      </Text>
      {typeof children === 'string' ? (
        <Text color={colors.dimText}>{children}</Text>
      ) : (
        children
      )}
    </Box>
  );
}
