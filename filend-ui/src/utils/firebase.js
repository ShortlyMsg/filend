import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";
import firebaseConfig from "@/config/firebaseConfig.json";

const firebaseApp = initializeApp(firebaseConfig);
const messaging = getMessaging(firebaseApp);

//console.log("Firebase initialized:", firebaseApp);
//console.log("Messaging initialized:", messaging);


export { messaging, getToken, onMessage };