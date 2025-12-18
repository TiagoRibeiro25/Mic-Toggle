import { useEffect, useState } from "react";
import { GetHotkey, SetHotkey } from "../../wailsjs/go/main/App";

type ModifierKey = "Ctrl" | "Shift" | "Alt";

export default function HotkeyPicker() {
	const [recording, setRecording] = useState<boolean>(false);
	const [hotkey, setHotkey] = useState<string>("");

	// Load saved hotkey on mount
	useEffect(() => {
		GetHotkey().then(setHotkey);
	}, []);

	function handleKeyDown(e: React.KeyboardEvent<HTMLDivElement>) {
		if (!recording) return;

		e.preventDefault();
		e.stopPropagation();

		const keys: ModifierKey[] = [];

		if (e.ctrlKey) keys.push("Ctrl");
		if (e.shiftKey) keys.push("Shift");
		if (e.altKey) keys.push("Alt");

		// Ignore pure modifier presses
		if (["Control", "Shift", "Alt"].includes(e.key)) {
			return;
		}

		const mainKey = e.key.length === 1
			? e.key.toUpperCase()
			: e.key;

		const combo = [...keys, mainKey].join("+");

		setHotkey(combo);
		setRecording(false);

		// Persist to Go backend
		SetHotkey(combo).catch(console.error);
	}

	return (
		<div
			tabIndex={0}
			onKeyDown={handleKeyDown}
			className="max-w-md p-4 bg-zinc-800 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
		>
			<p className="mb-2 font-semibold">Toggle Hotkey</p>

			<button
				type="button"
				onClick={() => setRecording(true)}
				className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500 transition"
			>
				{recording ? "Press keys…" : hotkey || "Set Hotkey"}
			</button>

			{recording && (
				<p className="mt-2 text-sm text-zinc-400">
					Press a key combination (e.g. Ctrl + Shift + M)
				</p>
			)}
		</div>
	);
}
