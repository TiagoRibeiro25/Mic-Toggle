//go:build windows
// +build windows

package mic

import (
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

// ToggleMic toggles the default microphone mute state and returns the new state.
func ToggleMic() (muted bool, err error) {
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	// Create device enumerator
	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	// Get default audio capture device (microphone)
	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	// Activate IAudioEndpointVolume
	var aev *wca.IAudioEndpointVolume
	if err = mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return
	}
	defer aev.Release()

	// Get current mute state
	if err = aev.GetMute(&muted); err != nil {
		return
	}

	// Toggle mute
	if err = aev.SetMute(!muted, nil); err != nil {
		return
	}

	return !muted, nil
}

func IsMuted() (bool, error) {
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return false, err
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return false, err
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.EConsole, &mmd); err != nil {
		return false, err
	}
	defer mmd.Release()

	var aev *wca.IAudioEndpointVolume
	if err := mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return false, err
	}
	defer aev.Release()

	var muted bool
	if err := aev.GetMute(&muted); err != nil {
		return false, err
	}

	return muted, nil
}
