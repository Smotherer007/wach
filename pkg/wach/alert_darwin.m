//go:build darwin
// +build darwin

#import <CoreFoundation/CoreFoundation.h>
#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>
#import <string.h>

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

// Show alert with an "Open Settings" button that opens Accessibility prefs.
// Returns 1 if the user clicked "Settings öffnen", 0 for "OK".
int showAlertWithSettingsLink(const char* title, const char* msg) {
	CFStringRef t = CFStringCreateWithCString(NULL, title, kCFStringEncodingUTF8);
	CFStringRef m = CFStringCreateWithCString(NULL, msg, kCFStringEncodingUTF8);

	CFOptionFlags result;
	CFUserNotificationDisplayAlert(
		0, kCFUserNotificationNoteAlertLevel,
		NULL, NULL, NULL,
		t, m,
		CFSTR("Einstellungen öffnen"),  // default button
		CFSTR("OK"),                     // alternate button
		NULL, &result);

	CFRelease(t);
	CFRelease(m);

	// kCFUserNotificationDefaultResponse = button 1 (default)
	// kCFUserNotificationAlternateResponse = button 2 (alternate)
	return (result == kCFUserNotificationDefaultResponse) ? 1 : 0;
}

void openURL(const char* url) {
	@autoreleasepool {
		NSString* nsUrl = [NSString stringWithUTF8String:url];
		[[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:nsUrl]];
	}
}

// Open Accessibility settings directly (macOS Ventura+)
void openAccessibilitySettings() {
	@autoreleasepool {
		NSURL *url = [NSURL URLWithString:@"x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"];
		[[NSWorkspace sharedWorkspace] openURL:url];
	}
}

// Check if the current process actually has Accessibility permissions
int hasAccessibilityPermission() {
	return AXIsProcessTrusted() ? 1 : 0;
}

int isDarkMode() {
	@autoreleasepool {
		NSString *style = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
		return [style isEqualToString:@"Dark"] ? 1 : 0;
	}
}

const char* systemLanguage() {
	@autoreleasepool {
		NSString *lang = [[NSLocale preferredLanguages] firstObject];
		if (lang == nil) return "en";
		// Extract language code (e.g. "de-DE" -> "de")
		NSArray *parts = [lang componentsSeparatedByString:@"-"];
		NSString *code = [parts firstObject];
		if (code == nil) return "en";
		return strdup([code UTF8String]);
	}
}
