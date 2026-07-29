//go:build darwin
// +build darwin

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

void showAlert(const char* title, const char* msg) {
	@autoreleasepool {
		NSAlert* alert = [[NSAlert alloc] init];
		alert.messageText = [NSString stringWithUTF8String:title];
		alert.informativeText = [NSString stringWithUTF8String:msg];
		[alert addButtonWithTitle:@"OK"];
		[alert runModal];
	}
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
