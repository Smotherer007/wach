package wach

/*
#cgo LDFLAGS: -framework CoreGraphics -framework Cocoa -framework IOKit

#include <CoreGraphics/CoreGraphics.h>
#include <IOKit/ps/IOPowerSources.h>
#include <IOKit/ps/IOPSKeys.h>
#include <IOKit/pwr_mgt/IOPMLib.h>
#include <stdlib.h>

double getIdleSeconds() {
	return (double)CGEventSourceSecondsSinceLastEventType(
		kCGEventSourceStateCombinedSessionState,
		kCGAnyInputEventType);
}

int tryMoveMouse(int dx, int dy) {
	CGEventRef event = CGEventCreate(NULL);
	CGPoint pos = CGEventGetLocation(event);
	CFRelease(event);

	CGPoint newPos = CGPointMake(pos.x + dx, pos.y + dy);
	CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
	CGEventRef move = CGEventCreateMouseEvent(source, kCGEventMouseMoved, newPos, kCGMouseButtonLeft);

	// Set delta values so the system recognizes this as real user activity
	CGEventSetIntegerValueField(move, kCGMouseEventDeltaX, dx);
	CGEventSetIntegerValueField(move, kCGMouseEventDeltaY, dy);

	CGEventPost(kCGHIDEventTap, move);
	CFRelease(move);
	CFRelease(source);

	CGEventRef check = CGEventCreate(NULL);
	CGPoint checkPos = CGEventGetLocation(check);
	CFRelease(check);

	return (checkPos.x == pos.x && checkPos.y == pos.y) ? 0 : 1;
}

int isDisplayAsleep() {
	CGDirectDisplayID display = CGMainDisplayID();
	if (display == kCGNullDirectDisplay) return 1;
	CGDisplayModeRef mode = CGDisplayCopyDisplayMode(display);
	if (mode == NULL) return 1;
	CFRelease(mode);
	return 0;
}

int getBatteryPercent() {
	CFTypeRef powerInfo = IOPSCopyPowerSourcesInfo();
	if (powerInfo == NULL) return -1;

	CFArrayRef list = IOPSCopyPowerSourcesList(powerInfo);
	if (list == NULL) {
		CFRelease(powerInfo);
		return -1;
	}

	int pct = -1;
	CFIndex count = CFArrayGetCount(list);
	if (count > 0) {
		CFDictionaryRef ps = IOPSGetPowerSourceDescription(powerInfo, CFArrayGetValueAtIndex(list, 0));
		if (ps != NULL) {
			CFNumberRef capacity = CFDictionaryGetValue(ps, CFSTR(kIOPSMaxCapacityKey));
			CFNumberRef current = CFDictionaryGetValue(ps, CFSTR(kIOPSCurrentCapacityKey));
			if (capacity != NULL && current != NULL) {
				int maxCap = 0, curCap = 0;
				CFNumberGetValue(capacity, kCFNumberIntType, &maxCap);
				CFNumberGetValue(current, kCFNumberIntType, &curCap);
				if (maxCap > 0) {
					pct = (int)((double)curCap / (double)maxCap * 100.0);
				}
			}
		}
	}

	CFRelease(list);
	CFRelease(powerInfo);
	return pct;
}

// ObjC functions implemented in alert_darwin.m
void showAlert(const char* title, const char* msg);
int isDarkMode();
void openURL(const char* url);
const char* systemLanguage();

// Create a power assertion to prevent idle sleep.
// Returns the assertion ID (>0) on success, or 0 on failure.
IOPMAssertionID createWakeAssertion(const char* reason) {
	IOPMAssertionID assertionID = 0;
	CFStringRef cfReason = CFStringCreateWithCString(kCFAllocatorDefault, reason, kCFStringEncodingUTF8);
	IOReturn result = IOPMAssertionCreateWithName(
		kIOPMAssertionTypePreventUserIdleSystemSleep,
		kIOPMAssertionLevelOn,
		cfReason,
		&assertionID);
	CFRelease(cfReason);
	return (result == kIOReturnSuccess) ? assertionID : 0;
}

// Release a previously created power assertion.
void releaseWakeAssertion(IOPMAssertionID assertionID) {
	if (assertionID != 0) {
		IOPMAssertionRelease(assertionID);
	}
}
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

const (
	AppName    = "Wach"
	AppVersion = "1.0.0"
)

var (
	AppAuthor  = "Patrick Weppelmann"
	AppSource  = "github.com/Smotherer007/wach"
	AppBasedOn = "github.com/prashantgupta24/automatic-mouse-mover"
)

func getIdleDuration() time.Duration {
	return time.Duration(float64(C.getIdleSeconds()) * float64(time.Second))
}

func tryMoveMouse(dx int) bool {
	return C.tryMoveMouse(C.int(dx), C.int(dx)) == 1
}

func isDisplayAsleep() bool {
	return C.isDisplayAsleep() == 1
}

// getBatteryPercent returns battery charge percentage (0-100) or -1 if not on battery.
func getBatteryPercent() int {
	return int(C.getBatteryPercent())
}

func IsDarkMode() bool {
	return C.isDarkMode() == 1
}

// SystemLanguage returns the user's preferred language code (e.g. "de", "en", "fr").
func SystemLanguage() string {
	cstr := C.systemLanguage()
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

func showAlert(title, msg string) {
	cTitle := C.CString(title)
	cMsg := C.CString(msg)
	C.showAlert(cTitle, cMsg)
	C.free(unsafe.Pointer(cTitle))
	C.free(unsafe.Pointer(cMsg))
}

// OpenGitHub opens the project page in the default browser.
func OpenGitHub() {
	url := C.CString("https://github.com/Smotherer007/wach")
	C.openURL(url)
	C.free(unsafe.Pointer(url))
}

// createWakeAssertion prevents the system from idle-sleeping.
// Returns the assertion ID (>0) on success, or 0 on failure.
func createWakeAssertion(reason string) uint32 {
	cReason := C.CString(reason)
	defer C.free(unsafe.Pointer(cReason))
	return uint32(C.createWakeAssertion(cReason))
}

// releaseWakeAssertion releases a previously created wake assertion.
func releaseWakeAssertion(id uint32) {
	C.releaseWakeAssertion(C.IOPMAssertionID(id))
}

// ShowAbout displays the About dialog with version info.
func ShowAbout() {
	l := global.locale
	msg := fmt.Sprintf(l.AboutMsg, AppVersion, l.AboutAuthor, l.AboutSource)
	showAlert(l.AboutTitle, msg)
}
