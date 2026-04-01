import React from 'react';
import { Box, Text } from 'ink';
import { colors } from '../theme.js';

// ASCII art — fits within 52 columns
const LOGO_LINES = [
  ' ██   ██  █████  ███   ██ ██ ███   ███  ██████ ',
  ' ██   ██ ██   ██ ████  ██ ██ ████ ████ ██    ██',
  ' ███████ ███████ ██ ██ ██ ██ ██ ███ ██ ██    ██',
  ' ██   ██ ██   ██ ██  ████ ██ ██     ██ ██    ██',
  ' ██   ██ ██   ██ ██   ███ ██ ██     ██  ██████ ',
];

// Compact version for narrow terminals (< 60 cols)
const LOGO_COMPACT = [
  ' ╻ ╻┏━┓┏┓╻╻┏┳┓┏━┓',
  ' ┣━┫┣━┫┃┗┫┃┃┃┃┃ ┃',
  ' ╹ ╹╹ ╹╹ ╹╹╹ ╹┗━┛',
];

// Gradient colors for each line of the banner
const GRADIENT = ['#89b4fa', '#89dceb', '#a6e3a1', '#f9e2af', '#fab387'];
const GRADIENT_COMPACT = ['#89b4fa', '#a6e3a1', '#fab387'];

interface BannerProps {
  version?: string;
  tagline?: string;
  cols?: number;
}

export const Banner = React.memo(function Banner({
  version = '0.1.0',
  tagline,
  cols = 80,
}: BannerProps): React.ReactElement {
  const useCompact = cols < 60;
  const lines = useCompact ? LOGO_COMPACT : LOGO_LINES;
  const gradient = useCompact ? GRADIENT_COMPACT : GRADIENT;

  return (
    <Box flexDirection="column" alignItems="center" width="100%" paddingTop={1}>
      {lines.map((line, i) => (
        <Box key={i} justifyContent="center" width="100%">
          <Text color={gradient[i % gradient.length]} bold>{line}</Text>
        </Box>
      ))}
      <Box justifyContent="center" width="100%" marginTop={1}>
        <Text color={colors.dimText}>
          {'v' + version}
          {tagline ? ` \u2014 ${tagline}` : ''}
        </Text>
      </Box>
    </Box>
  );
});
