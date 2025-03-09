import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";
import firebaseConfig from "@/config/firebaseConfig.json";

let messaging = null
export function initializeFirebase(){
    const firebaseApp = initializeApp(firebaseConfig);
    messaging = getMessaging(firebaseApp);

    console.log("Firebase initialized:", firebaseApp);
    console.log("Messaging initialized:", messaging);
}


export { messaging, getToken, onMessage };