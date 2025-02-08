let isError = false

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function displayError(err) {
    // do not display error if there is an error already set
    if(!isError) {
        // stop animation
        document.getElementById("sheets").classList.add("no-animation");

        // set error message
        document.getElementById("error").style.display = "block";
        document.getElementById("error").innerHTML = err;

        isError = true
    }
}

async function uploadStatus(timeoutCalculated) {
    while(Date.now() < timeoutCalculated) {                
        try {
            const response = await fetch('/status')
            const data = await response.json()

            console.log(data)

            if(data.uploadFinished === true) {
                location.replace(data.webLink);
                break
            } else {
                if(data.hasOwnProperty('uploadError')) {
                    displayError(data.uploadError)
                    break
                }
            }
        } catch(err) {
            console.log(err)
        }

        await sleep(1000);
    }

    displayError("Timed out. Server did not respond in time.");
}

// keep alive for 
timeout = 20;
const timeoutCalculated = Date.now() + (timeout * 1000) 

// auth
const fragment = new URLSearchParams(window.location.hash.slice(1));
console.log("Fragment:", fragment);
if(fragment.size > 0) {
    const accessToken = fragment.get('access_token');
    const tokenType = fragment.get('token_type');
    const expiresIn = parseInt(fragment.get('expires_in'), 10);
    const state = fragment.get('state');
    const scope = fragment.get('scope');
    const error = fragment.get('error');

    // TODO handle error

    fetch('/auth', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({accessToken: accessToken, tokenType: tokenType, expiresIn: expiresIn, state: state})
    })
}

// check on upload status
uploadStatus(timeoutCalculated);
