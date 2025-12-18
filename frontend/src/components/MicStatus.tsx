import { useEffect, useState } from "react";
import { GetMicState } from "../../wailsjs/go/main/App";
import { EventsOn } from "../../wailsjs/runtime";

export default function MicStatus() {
  const [muted, setMuted] = useState<boolean>(false);

  // Load initial state
  useEffect(() => {
    GetMicState().then(setMuted);
  }, []);

  // Listen for backend events
  useEffect(() => {
    const unsubscribe = EventsOn("micStateChanged", (state: boolean) => {
      setMuted(state);
    });
    return () => unsubscribe();
  }, []);

  return (
    <div className="mt-4 p-4 bg-zinc-800 rounded-xl max-w-md">
      <p className="font-semibold">Microphone Status:</p>
      <p className={muted ? "text-red-500" : "text-green-500"}>
        {muted ? "Muted" : "Unmuted"}
      </p>
    </div>
  );
}
