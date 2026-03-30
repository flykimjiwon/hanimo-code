import { execaCommand } from 'execa';

/**
 * Send a desktop notification when a long-running task completes.
 * Supports macOS (osascript), Linux (notify-send), and Windows (PowerShell).
 * Falls back silently on unsupported platforms or missing tools.
 */
export async function sendNotification(
  title: string,
  message: string,
): Promise<void> {
  const clean = (s: string): string => s.replace(/["'\\`$]/g, '').slice(0, 200);
  const safeTitle = clean(title);
  const safeMessage = clean(message);

  try {
    if (process.platform === 'darwin') {
      await execaCommand(
        `osascript -e 'display notification "${safeMessage}" with title "${safeTitle}"'`,
        { shell: true, timeout: 5000, reject: false },
      );
    } else if (process.platform === 'linux') {
      await execaCommand(
        `notify-send "${safeTitle}" "${safeMessage}"`,
        { shell: true, timeout: 5000, reject: false },
      );
    } else if (process.platform === 'win32') {
      const ps = `powershell -Command "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('${safeMessage}','${safeTitle}')"`;
      await execaCommand(ps, { shell: true, timeout: 5000, reject: false });
    }
  } catch {
    // Silent fail — notification is non-critical
  }
}

/**
 * Play a terminal bell as a simple audible notification.
 */
export function bell(): void {
  process.stdout.write('\x07');
}
