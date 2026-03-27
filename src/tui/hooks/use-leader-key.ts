import { useState, useCallback, useRef, useEffect } from 'react';

export type LeaderAction =
  | 'palette'
  | 'save'
  | 'load'
  | 'verbose'
  | 'clear'
  | 'help'
  | 'sessions';

const KEY_MAP: Record<string, LeaderAction> = {
  k: 'palette',
  s: 'save',
  l: 'load',
  v: 'verbose',
  c: 'clear',
  h: 'help',
  e: 'sessions',
};

const LEADER_TIMEOUT_MS = 2000;

export interface UseLeaderKeyReturn {
  leaderActive: boolean;
  processKey: (input: string, key: { ctrl: boolean; escape: boolean }) => LeaderAction | null;
}

export function useLeaderKey(): UseLeaderKeyReturn {
  const [leaderActive, setLeaderActive] = useState(false);
  const activeRef = useRef(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  useEffect(() => () => clearTimer(), [clearTimer]);

  const processKey = useCallback(
    (input: string, key: { ctrl: boolean; escape: boolean }): LeaderAction | null => {
      // Second key in leader mode
      if (activeRef.current) {
        activeRef.current = false;
        setLeaderActive(false);
        clearTimer();

        // Ctrl or Escape cancels leader — let caller handle normally
        if (key.ctrl || key.escape) return null;

        return KEY_MAP[input.toLowerCase()] ?? null;
      }

      // Ctrl+X activates leader mode
      if (key.ctrl && input === 'x') {
        activeRef.current = true;
        setLeaderActive(true);
        clearTimer();
        timerRef.current = setTimeout(() => {
          activeRef.current = false;
          setLeaderActive(false);
        }, LEADER_TIMEOUT_MS);
        return null;
      }

      return null;
    },
    [clearTimer],
  );

  return { leaderActive, processKey };
}
