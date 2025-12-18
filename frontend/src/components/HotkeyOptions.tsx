import { useEffect, useState } from "react";
import {
  GetPlayBeep,
  GetShowNotification,
  SetPlayBeep,
  SetShowNotification,
} from "../../wailsjs/go/main/App";

export default function HotkeyOptions() {
	const [playBeep, setPlayBeep] = useState(false);
	const [showNotification, setShowNotification] = useState(false);

	// Load initial values from backend
	useEffect(() => {
		GetPlayBeep().then(setPlayBeep);
		GetShowNotification().then(setShowNotification);
	}, []);

	const handlePlayBeepChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
		const checked = e.target.checked;
		setPlayBeep(checked);
		await SetPlayBeep(checked);
	};

	const handleShowNotificationChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
		const checked = e.target.checked;
		setShowNotification(checked);
		await SetShowNotification(checked);
	};

	return (
		<div className="mt-6 space-y-4">
			<label className="flex items-center space-x-2">
				<input
					type="checkbox"
					checked={playBeep}
					onChange={handlePlayBeepChange}
					className="form-checkbox"
				/>
				<span>Play beep on hotkey</span>
			</label>

			<label className="flex items-center space-x-2">
				<input
					type="checkbox"
					checked={showNotification}
					onChange={handleShowNotificationChange}
					className="form-checkbox"
				/>
				<span>Show notification on hotkey</span>
			</label>
		</div>
	);
}
