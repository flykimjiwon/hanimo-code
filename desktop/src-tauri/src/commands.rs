use std::io::BufRead;
use tauri::{AppHandle, Emitter, Manager, State};

use crate::sidecar::Sidecar;

#[tauri::command]
pub async fn start_sidecar(
    app: AppHandle,
    sidecar: State<'_, Sidecar>,
    provider: String,
    model: String,
    api_key: String,
    role: Option<String>,
    base_url: Option<String>,
    project_dir: Option<String>,
) -> Result<String, String> {
    let node_path = std::env::var("NODE_PATH_OVERRIDE").unwrap_or_else(|_| "node".to_string());

    let script_path = if cfg!(debug_assertions) {
        "../../dist/cli.js".to_string()
    } else {
        let resource_dir = app
            .path()
            .resource_dir()
            .map_err(|e| format!("Failed to get resource dir: {}", e))?;
        resource_dir
            .join("dist")
            .join("cli.js")
            .to_string_lossy()
            .to_string()
    };

    let mut args = vec!["--headless".to_string()];
    args.push("-p".to_string());
    args.push(provider.clone());
    args.push("-m".to_string());
    args.push(model.clone());
    args.push("-k".to_string());
    args.push(api_key.clone());
    if let Some(ref r) = role {
        args.push("--role".to_string());
        args.push(r.clone());
    }
    if let Some(ref u) = base_url {
        args.push("-u".to_string());
        args.push(u.clone());
    }

    let env_vars: &[(&str, &str)] = &[("HANIMO_FF_HEADLESS_MODE", "1")];

    sidecar.start(&node_path, &script_path, &args, env_vars, project_dir.as_deref())?;

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

    if let Some(stderr) = sidecar.take_stderr() {
        let app2 = app.clone();
        std::thread::spawn(move || {
            let reader = std::io::BufReader::new(stderr);
            for line in reader.lines() {
                if let Ok(text) = line {
                    let _ = app2.emit("sidecar-stderr", text);
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
