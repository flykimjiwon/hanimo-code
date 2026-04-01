use std::io::Write;
use std::process::{Child, ChildStdout, Command, Stdio};
use std::sync::Mutex;

pub struct Sidecar {
    process: Mutex<Option<Child>>,
    stdout: Mutex<Option<ChildStdout>>,
}

impl Sidecar {
    pub fn new() -> Self {
        Sidecar {
            process: Mutex::new(None),
            stdout: Mutex::new(None),
        }
    }

    pub fn start(&self, node_path: &str, script_path: &str, args: &[&str]) -> Result<(), String> {
        let mut child = Command::new(node_path)
            .arg(script_path)
            .args(args)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .stderr(Stdio::null())
            .spawn()
            .map_err(|e| format!("Failed to spawn sidecar: {}", e))?;

        let stdout = child
            .stdout
            .take()
            .ok_or_else(|| "Failed to capture stdout".to_string())?;

        *self.process.lock().unwrap() = Some(child);
        *self.stdout.lock().unwrap() = Some(stdout);

        Ok(())
    }

    pub fn send(&self, message: &str) -> Result<(), String> {
        let mut guard = self.process.lock().unwrap();
        if let Some(child) = guard.as_mut() {
            if let Some(stdin) = child.stdin.as_mut() {
                writeln!(stdin, "{}", message)
                    .map_err(|e| format!("Failed to write to stdin: {}", e))?;
            } else {
                return Err("Sidecar stdin not available".to_string());
            }
        } else {
            return Err("Sidecar not running".to_string());
        }
        Ok(())
    }

    pub fn take_stdout(&self) -> Option<ChildStdout> {
        self.stdout.lock().unwrap().take()
    }

    pub fn stop(&self) -> Result<(), String> {
        let mut guard = self.process.lock().unwrap();
        if let Some(mut child) = guard.take() {
            child.kill().map_err(|e| format!("Failed to kill sidecar: {}", e))?;
        }
        Ok(())
    }
}
