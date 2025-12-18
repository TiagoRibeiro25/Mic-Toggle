import HotkeyOptions from "./components/HotkeyOptions";
import HotkeyPicker from "./components/HotkeyPicker";
import MicStatus from "./components/MicStatus";

function App() {
	return (
		<div className="min-h-screen bg-zinc-900 text-zinc-100 p-6">
			<h1 className="text-2xl font-bold mb-4">Mic Toggle Settings</h1>
			
			<HotkeyPicker />
			<HotkeyOptions />
			<MicStatus />
		</div>
	);
}

export default App;
