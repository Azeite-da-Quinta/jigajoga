import {useContext, useState} from "react";
import {AppContext} from "../AppContext/AppContext.tsx";
import styles from "./ConnectionForm.module.css"

const ConnectionForm = () => {

    const { handleUrl, connState, wsFuncs, reconnect, handleReconnect, disableIfConnected, disableIfDisconnected } = useContext(AppContext);

    const generateUrlString = (event) => {
        let children = event.target.children;
        return children[1].value + children[2].value;
    }

    const formAction = (event) => {
        event.preventDefault();
        let currValue = generateUrlString(event);
        if(event.nativeEvent.submitter.name == "connect"){
            console.info("ConnectionForm: new connection to " + currValue);
            handleUrl(currValue);
            handleReconnect(reconnect => reconnect + 1 )
        } else if (event.nativeEvent.submitter.name == "disconnect") {
            console.info("ConnectionForm: requesting disconnect from: "+ currValue);
            wsFuncs.disconnect();
        }
    }

    const connStatus = () => {
        let ret = "Not Connected";
        switch (connState) {
            case 0:
                ret = "Connecting...";
                break;
            case 1:
                ret = "Connected!";
                break;
            case 2:
                ret = "Disconnecting...";
                break;
            case 3:
                ret = "Disconnected!";
                break;
            case 4:
                ret = "Couldn't Connect!"
        }
        return ret;
    }



    return (
        <nav>
            <span>
                Jigajoga web-based client debugger v1
            </span>
            <span>
                <form onSubmit={formAction}>
                    <label>Endpoint: </label>
                    <select disabled={disableIfConnected()}>
                        <option value="ws://">ws://</option>
                        <option value="wss://">wss://</option>
                    </select>
                    <input type="text" placeholder={"address.domain:port/path"} disabled={disableIfConnected()} />
                    <input type="submit" name="connect" value="Connect" disabled={disableIfConnected()}/>
                    <input type="submit" name="disconnect" value="Disconnect" disabled={disableIfDisconnected()}/>
                </form>
            </span>
            <span>

            </span>
            <span>
                {connStatus()}
            </span>
        </nav>
    );
}

export default ConnectionForm