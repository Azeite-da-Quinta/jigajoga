import {useContext, useState} from "react";

const RequestForm = ( { sendMsg, disabled  } ) => {

        const [reqMessage, handleReqMessage] = useState("{\"your\":\"message\"}");

        const changeMessage = (event) => {
            handleReqMessage(() => event.target.value)
        }

        const submitMessage = (event)  => {
            event.preventDefault();
            sendMsg(reqMessage);
        }

        return (

                <form onSubmit={submitMessage} >
                    <label>Endpoint</label>
                    <input type="text" placeholder={"{\"your\":\"message\"}"}  onChange={changeMessage} disabled={disabled()}/>
                    <input type="submit" value="Submit" disabled={disabled()}/>
                </form>

        );
    }


export default RequestForm;