use std::io::BufRead;
use tauri::{AppHandle, Emitter, State};

use crate::sidecar::Sidecar;

#[tauri::command]
pub fn start_sidecar(app: AppHandle, sidecar: State<Sidecar>) -> Result<String, String> {
    let node_path = "node";
    let script_path = "../../dist/headless.js";

    sidecar.start(node_path, script_path, &["--headless"])?;

    if let Some(stdout) = sidecar.take_stdout() {
        let app_clone = app.clone();
        std::thread::spawn(move || {
            let reader = std::io::BufReader::new(stdout);
            for line in reader.lines() {
                match line {
                    Ok(text) => {
                        let _ = app_clone.emit("sidecar-event", text);
                    }
                    Err(_) => break,
                }
            }
        });
    }

    Ok("Sidecar started".to_string())
}

#[tauri::command]
pub fn send_prompt(sidecar: State<Sidecar>, content: String) -> Result<(), String> {
    let message = serde_json::json!({
        "type": "prompt",
        "content": content
    });
    sidecar.send(&message.to_string())
}

#[tauri::command]
pub fn stop_sidecar(sidecar: State<Sidecar>) -> Result<(), String> {
    sidecar.stop()
}
