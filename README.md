![](extras/art/guppy-1.png)

# guppy
Guppy is a simple tool for opening files from your computer in Google Drive. It can be set as the default file handler 
or as a command-line tool.

## Usage
When you launch guppy for the first time, it will set itself as the default handler. Then you can double-click or 
right-click and select open with any CSV, Excel, Word, or PowerPoint file to have guppy upload and open it in Google 
Drive. 

#### Command Line 
`guppy [file]` `guppy upload [file]`  
Upload a file to Google Drive.  

`guppy login`  
Login to Google Drive.  

`guppy help [command]`  
Display help information.

#### Flags
```
-h, --help      Provides help command
-v, --verbose   Enable verbose output
```

## Installation
To use guppy, you will need to set up a Google OAuth2 Client ID. As an organization, you can distribute the OAuth2 
Client ID to your users via an MDM configuration profile.

#### Create a Google OAuth2 Client ID for Desktop applications
1. Go to https://console.cloud.google.com/projectcreate
2. Enable the Google Drive API
3. Configure OAuth consent screen
4. Create a OAuth2 Client ID for Web application, save the client ID
5. Add `http://localhost/interstitial.html` as a redirect URI

#### Deploy a configuration profile for macOS
Examples of configuration profiles can be found in the `extras` folder. Substitute the client ID in the configuration 
profile with your own from Step 4 above.
| key                  | type    | required | default |
|----------------------|---------|----------|---------|
| GoogleOauth2ClientId | String  | true     |         |
| DisableAnalytics     | Boolean | false    | false   |

#### Download the guppy app from releases
The latest release of guppy can be found in the [releases](https://github.com/forevanyeung/guppy/releases) section of 
this repository.
