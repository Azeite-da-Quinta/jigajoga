interface Props {
    endpoint: string;
    request: JSON;
    response: JSON;
}

const log = ( {request,response} : Props) => {
    return(
        <>
            <p><a>Request:</a><a> {JSON.stringify(request) } </a></p>
            <p><a>Response:</a><a> {JSON.stringify(response) } </a></p>
        </>
    );
}

export default log;