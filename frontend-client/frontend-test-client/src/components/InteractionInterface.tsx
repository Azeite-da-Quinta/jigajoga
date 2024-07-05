import RequestForm from "./Forms/RequestForm"
import RequestList from "./RequestsLog/RequestList.tsx";
import {useContext} from "react";
import {AppContext} from "./AppContext/AppContext.tsx";

const InteractionInterface = () => {

    const { wsFuncs, disableIfDisconnected } = useContext(AppContext);


    return(<div>
            <section>
                <RequestForm sendMsg={wsFuncs.sendMsg} disabled={disableIfDisconnected} />
            </section>
            <section>
                <RequestList/>
            </section>
        </div>)

}

export default InteractionInterface;