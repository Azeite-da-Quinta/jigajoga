import {useEffect, createContext ,useState} from "react";

export const AppContext : React.Context<object> = createContext(
    {
        url: "",
        handleUrl: () => {},
        disableIfConnected: () => {},
        disableIfDisconnected: () => {},
        connState: null,
        wsFuncs: {
            sendMsg: (value) => {
                console.error("WebSocket wasn't initalized.")
            },
            disconnect: () => {
                console.error("WebSocket wasn't initalized.")
            }
        },
        reconnect: 0,
        handleReconnect: () => {},
        logs: [],
        request: null,
    }
);

const AppContextProvider = ({children} : Props) => {

    const [url, handleUrl] = useState("");
    const [connState, handleConnState] = useState(-1);
    const [reconnect, handleReconnect] = useState(0);
    const [logs, handleLogs] = useState([]);
    const [wsFuncs, handleFuncs] = useState({
        sendMsg: (value) => {
            console.error("WebSocket wasn't initalized.");
        },
        disconnect: () => {
            console.error("WebSocket wasn't initalized.");
        }
    })
    const [request, handleRequest] = useState("");
    const [response, handleResponse] = useState("");

    const pushLog = (log) => {
        handleLogs(logs => [...logs,log])
    }

    /*
     let disconnect = () => {
        /* Mocked Behaviour for interface
        handleConnState(2)
        setTimeout(() => {
                console.log("Disconnected!")
                handleConnState(3)
            },
            5000)
        console.info("Closing connection with " + ws.url)
        handleConnState(2);
        ws.close();
    }*/

    useEffect(() => {
        if(request != "" && response != ""){
            pushLog({request:request,response:response})
            handleRequest("")
            handleResponse("")
        }
    }, [response]);

    useEffect(() => {

        if(url != ""){
            /* Mocked Behaviour for interface
            console.info("changed endpoint to:" + url)
            handleConnState(0)
            setTimeout(() => {
                console.log("Connected!")
                handleConnState(1)
                },
                5000) */
            handleConnState(0)
            const ws = new WebSocket(
                url,
                [
                    'base64url.bearer.authorization.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1Ijp7InIiOiIxMjM0IiwiaWQiOiIxIiwibiI6ImJvYiJ9LCJpc3MiOiJqaWdham9nYSIsImV4cCI6MTcyMDIyMTQ0OH0.8vdhWY8tljQ7DIXQL9DCKiUUpXZ3mX49A5wzTn6r9xU',
                    "v0.jigajoga.json"]);

            handleFuncs({
                sendMsg: (value) => {
                    handleRequest(() => value);
                    ws.send(value);
                },
                disconnect: () => {
                    console.info("Closing connection with " + ws.url);
                    handleConnState(2);
                    ws.close();
                }
            });

            ws.addEventListener("open", (event : Event) => {
                console.log("Connecting to: " + ws.url);
                handleConnState(1)
            })

            ws.addEventListener("message", (event) => {
                console.log(event);
                handleResponse(event.data);
            });

            ws.addEventListener("close", (event) => {
                console.log("Closing connection to " + ws.url);
                handleConnState(3)
            })

            ws.addEventListener("error", (event) => {
                console.error("Couldn't connect at this time.");
                ws.close();
                handleConnState(4)
            })
        }

        return (ws:WebSocket) => {
            if(ws != undefined) {
                ws.close();
            }
        }

    }, [reconnect]);

    const disableIfConnected = () => {
        return connState != 3 && connState != -1;
    }

    const disableIfDisconnected = () =>  {
        return connState != 1;
    }

    return(
        <AppContext.Provider value={ {url, handleUrl, connState, disableIfConnected, disableIfDisconnected, wsFuncs, reconnect, handleReconnect,logs , request} }>
            {children}
        </AppContext.Provider>
    )
}

export default AppContextProvider;