package wach

/*
#cgo LDFLAGS: -framework CoreGraphics -framework Cocoa

#include <CoreGraphics/CoreGraphics.h>

// CoreGraphics-based idle detection and mouse movement (pure C, no ObjC)

double getIdleSeconds() {
	CFTimeInterval idle = CGEventSourceSecondsSinceLastEventType(
		kCGEventSourceStateCombinedSessionState,
		kCGAnyInputEventType
	);
	return (double)idle;
}

int tryMoveMouse(int dx, int dy) {
	CGEventRef event = CGEventCreate(NULL);
	CGPoint pos = CGEventGetLocation(event);
	CFRelease(event);

	CGPoint newPos = CGPointMake(pos.x + dx, pos.y + dy);
	CGEventRef move = CGEventCreateMouseEvent(NULL, kCGEventMouseMoved, newPos, kCGMouseButtonLeft);
	CGEventPost(kCGHIDEventTap, move);
	CFRelease(move);

	CGEventRef check = CGEventCreate(NULL);
	CGPoint checkPos = CGEventGetLocation(check);
	CFRelease(check);

	if (checkPos.x == pos.x && checkPos.y == pos.y) {
		return 0;
	}
	return 1;
}

int isDisplayAsleep() {
	CGDirectDisplayID display = CGMainDisplayID();
	if (display == kCGNullDirectDisplay) return 1;
	CGDisplayModeRef mode = CGDisplayCopyDisplayMode(display);
	if (mode == NULL) return 1;
	CFRelease(mode);
	return 0;
}

// Implemented in alert_darwin.m (Objective-C)
void showAlert(const char* title, const char* msg);
int isDarkMode();
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

// App metadata
const (
	AppName    = "Wach"
	AppVersion = "1.0.0"
)

var (
	AppAuthor  = "Patrick Weppelmann"
	AppSource  = "github.com/patrickweppelmann/wach"
	AppBasedOn = "github.com/prashantgupta24/automatic-mouse-mover"
)

// getIdleDuration returns how long the system has been idle.
func getIdleDuration() time.Duration {
	secs := float64(C.getIdleSeconds())
	return time.Duration(secs * float64(time.Second))
}

// tryMoveMouse attempts to move the mouse by (dx, dy).
func tryMoveMouse(dx int) bool {
	return C.tryMoveMouse(C.int(dx), C.int(dx)) == 1
}

// isDisplayAsleep returns true if the main display is sleeping.
func isDisplayAsleep() bool {
	return C.isDisplayAsleep() == 1
}

// showAlert shows a native macOS alert dialog via Objective-C.
func showAlert(title, msg string) {
	cTitle := C.CString(title)
	cMsg := C.CString(msg)
	C.showAlert(cTitle, cMsg)
	C.free(unsafe.Pointer(cTitle))
	C.free(unsafe.Pointer(cMsg))
}

// IsDarkMode returns true if macOS is in Dark Mode appearance.
func IsDarkMode() bool {
	return C.isDarkMode() == 1
}

// ShowAbout displays the About dialog.
func ShowAbout() {
	msg := fmt.Sprintf(
		"Version %s\n\n"+
			"Ein minimaler Mausbeweger fur macOS (Apple Silicon).\n"+
			"Halt den Mac wach, indem er bei Inaktivitat die Maus bewegt.\n\n"+
			"(c) %s\n%s\n\n"+
			"Basierend auf:\n%s",
		AppVersion, AppAuthor, AppSource, AppBasedOn,
	)
	showAlert("Uber Wach", msg)
}
