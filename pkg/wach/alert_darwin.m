//go:build darwin
// +build darwin

#import <CoreFoundation/CoreFoundation.h>
#import <Cocoa/Cocoa.h>

// Thread-safe alert using CoreFoundation (no main thread requirement)
void showAlert(const char* title, const char* msg) {
	CFStringRef t = CFStringCreateWithCString(NULL, title, kCFStringEncodingUTF8);
	CFStringRef m = CFStringCreateWithCString(NULL, msg, kCFStringEncodingUTF8);

	CFOptionFlags result;
	CFUserNotificationDisplayAlert(
		0, kCFUserNotificationNoteAlertLevel,
		NULL, NULL, NULL,
		t, m, CFSTR("OK"),
		NULL, NULL, &result);

	CFRelease(t);
	CFRelease(m);
}

void openURL(const char* url) {
	@autoreleasepool {
		NSString* nsUrl = [NSString stringWithUTF8String:url];
		[[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:nsUrl]];
	}
}

int isDarkMode() {
	@autoreleasepool {
		NSString *style = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
		return [style isEqualToString:@"Dark"] ? 1 : 0;
	}
}
