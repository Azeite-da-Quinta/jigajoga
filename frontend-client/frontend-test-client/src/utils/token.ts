import * as jose from 'jose';

// just for testing purposes
const secret = "Alheiras"; //"QWxoZWlyYXM=";

const alg = 'HS256'
const typ = "JWT"

const contents =  { "u" :{
    "r" : "1234",
    "id" : "1",
    "n" : "frontend-debugger"
} };

export const jwt = await new jose.SignJWT(contents)
    .setProtectedHeader({ alg, typ })
    .setIssuer("jigajoga")
//    .setIssuedAt()
    .setExpirationTime("2h")
    .sign(new TextEncoder().encode(secret));

console.log(jwt)

const auth_header = "base64url.bearer.authorization."+jwt

export default auth_header;