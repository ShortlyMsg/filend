import { messaging, getToken } from "@/firebaseConfig";

async function requestNotificationPermission() {
    try {
        const permission = await Notification.requestPermission();
        if (permission === "granted") {
            console.log("Bildirim izni verildi.");
            return await getFCMToken();
        } else {
            console.log("Bildirim izni verilmedi.");
            return null;
        }
    } catch (error) {
        console.error("Bildirim izni istenirken hata oluştu:", error);
        return null;
    }
}

async function getFCMToken() {
    try {
        const token = await getToken(messaging, { vapidKey: "YOUR_VAPID_KEY" }); // VAPID key
        if (token) {
            console.log("FCM token alındı:", token);
            return token;
        } else {
            console.log("FCM token alınamadı.");
            return null;
        }
    } catch (error) {
        console.error("FCM token alınırken hata oluştu:", error);
        return null;
    }
}

export { requestNotificationPermission, getFCMToken };