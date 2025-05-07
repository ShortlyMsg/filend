importScripts('https://www.gstatic.com/firebasejs/10.11.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.11.0/firebase-messaging-compat.js');

firebase.initializeApp({
    "apiKey": "api-key",
    "authDomain": "filend-msg.firebaseapp.com",
    "projectId": "filend-msg",
    "storageBucket": "filend-msg.firebasestorage.app",
    "messagingSenderId": "messaging-sender-id",
    "appId": "app-id",
});

const messaging = firebase.messaging();

messaging.onBackgroundMessage(function(payload) {
    console.log('[firebase-messaging-sw.js] Background message received:', payload);
    const notificationTitle = payload.notification.title;
    const notificationOptions = {
        body: payload.notification.body
    };

    self.registration.showNotification(notificationTitle, notificationOptions);
});