//
//  guppyApp.swift
//  guppy
//
//  Created by Evan Yeung on 2/12/25.
//

import SwiftUI
import Foundation

@main
struct guppyApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate

    var body: some Scene {
        Settings() {
            SettingsView()
        }
    }
}

func isFirstLaunch() -> Bool {
    let key = "hasLaunchedBefore"
    let launchedBefore = UserDefaults.standard.bool(forKey: key)

    if !launchedBefore {
        UserDefaults.standard.set(true, forKey: key)
        UserDefaults.standard.synchronize()
        return true
    }

    return false
}

func openSettings() {
    let environment = EnvironmentValues()
    environment.openSettings()
    NSApp.setActivationPolicy(.regular)
    NSApp.activate(ignoringOtherApps: true)
}

class AppDelegate: NSObject, NSApplicationDelegate {
    override init() {
        super.init()
    }
    
    func application(_ application: NSApplication, open urls: [URL]) {
        for url in urls {
            processFile(at: url)
        }
        
        #if !DEBUG
        // Terminate the app after processing
        NSApplication.shared.terminate(nil)
        #endif
    }
    
    func applicationDidFinishLaunching(_ notification: Notification) {
        if ProcessInfo.processInfo.environment["XCODE_RUNNING_FOR_PREVIEWS"] == "1" {
            return
        }
        
        if(isFirstLaunch()) {
            NSLog("Setting guppy as default on first launch")
            setDefaultHandler(for: contentTypes)
        }
        
        // Filter and process command line arguments
        let relevantArguments = CommandLine.arguments.dropFirst().filter { arg in
            !arg.hasPrefix("-") && arg != "YES" && !arg.hasSuffix(".app")
        }
        
        let urls = relevantArguments.map { URL(fileURLWithPath: $0) }
        
        if !urls.isEmpty {
            application(NSApplication.shared, open: urls)
        } else {
            // Open settings if app is opened directly without any file
            NSLog("guppy opened with no files, defaulting to opening settings")
            openSettings()
        }
    }
    
    private func processFile(at url: URL) {
        NSLog("File opened: \(url.path)")
        
        guard let guppyBin = locateBinary(name: "guppy") else {
            NSLog("Could not find guppy binary")
            return
        }
            
        let process = Process()
        process.executableURL = URL(fileURLWithPath: guppyBin)
        process.arguments = [url.path, "--desktop"]
        #if DEBUG
        process.arguments! += ["-v"]
        #endif
        
        NSLog("Launched guppy with arguments: \(process.arguments ?? [])")
                
        let pipe = Pipe()
        process.standardOutput = pipe
        process.standardError = pipe

        let fileHandle = pipe.fileHandleForReading

        fileHandle.readabilityHandler = { fileHandle in
            let data = fileHandle.availableData
            if let output = String(data: data, encoding: .utf8), !output.isEmpty {
                NSLog(output)
            }
        }

        do {
            try process.run()
            process.waitUntilExit()
            
            if process.terminationStatus == 0 {
                NSLog("Guppy processed the file successfully")
            } else {
                NSLog("Guppy encountered an error. Exit code: \(process.terminationStatus)")
            }
        } catch {
            print("Failed to run process: \(error)")
        }

        // Remove the readability handler when done
        fileHandle.readabilityHandler = nil
    }
    
    private func locateBinary(name: String) -> String? {
        guard let bundlePath = Bundle.main.bundlePath as NSString? else {
            return nil
        }
        
        let binaryPath = bundlePath.appendingPathComponent("Contents/Resources/\(name)")
        return FileManager.default.fileExists(atPath: binaryPath) ? binaryPath : nil
    }
}
