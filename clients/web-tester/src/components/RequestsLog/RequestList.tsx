import {useContext} from "react";
import {AppContext} from "../AppContext/AppContext.tsx";
import RequestLog from "./RequestLog.tsx";


const RequestList = () => {

    const { logs } = useContext(AppContext);

    return (
        <div>
            <p>Requests made to Endpoint</p>
            {logs.map((element: object, index: number) => {
                return(<RequestLog key={index} request={element.request} response={element.response}/> );
            })}
        </div>
    )
}

export default RequestList;