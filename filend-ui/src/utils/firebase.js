import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";
import firebaseConfig from "@/config/firebaseConfig.json";

const messaging = {};
export function initializeFirebase(){
    const firebaseApp = initializeApp(firebaseConfig);
    messaging.current = getMessaging(firebaseApp);

    console.log("Firebase initialized:", firebaseApp);
    console.log("Messaging initialized:", messaging.current);
}


export { messaging, getToken, onMessage };