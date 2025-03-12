import SwiftUI
import Cocoa
import UniformTypeIdentifiers

let contentTypes = [
    UTType.commaSeparatedText,
    UTType.spreadsheet,
    UTType("org.openxmlformats.wordprocessingml.document"),
    UTType.presentation
].compactMap{$0} // removes nil

struct SettingsView: View {
    @State private var isDefault: Bool = false
    let timer = Timer.publish(every: 1, on: .main, in: .common).autoconnect()
    
    var body: some View {
        Form {
            LabeledContent("Version:", value: Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String ?? "0.0.0")
                .padding(1)
            
            LabeledContent("Default:") {
                VStack(alignment: .leading) {
                    if(isDefault) {
                        Text("guppy is the default app")
                    } else {
                        Text("guppy is not the default app")
                    }
                    Button("Set as default") {
                        setDefaultHandler(for: contentTypes)
                    }.disabled(isDefault)
                }
            }
        }
        .padding()
        .frame(width: 300, height: 200)
        .onAppear {
            // Check initial state
            isDefault = isDefaultHandler(for: contentTypes)
        }
        .onReceive(timer) { _ in
            // Periodically check if status has changed
            isDefault = isDefaultHandler(for: contentTypes)
        }
    }
}

#Preview {
    SettingsView()
}

func isDefaultHandler(for contentTypes: [UTType]) -> Bool {
    for contentType in contentTypes {
        guard let defaultAppURL = NSWorkspace.shared.urlForApplication(toOpen: contentType) else {
            return false // type does not have a default set
        }
        let defaultAppBundleID = Bundle(url: defaultAppURL)?.bundleIdentifier
        
        guard defaultAppBundleID == Bundle.main.bundleIdentifier else {
            return false // If any file type does NOT match, return false
        }
    }
    return true // Only returns true if ALL content types match
}

func setDefaultHandler(for contentTypes: [UTType]) {
    let appURL = Bundle.main.bundleURL
    
    for contentType in contentTypes {
        NSWorkspace.shared.setDefaultApplication(at: appURL, toOpen: contentType) { error in
            if let error = error {
                NSLog("Failed to set default handler: \(error.localizedDescription)")
            } else {
                NSLog("Default application request sent successfully. \(contentType.identifier)")
            }
        }
    }
}
