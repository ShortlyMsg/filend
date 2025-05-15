importScripts('https://www.gstatic.com/firebasejs/10.11.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.11.0/firebase-messaging-compat.js');

firebase.initializeApp({
    "apiKey": "AIzaSyAuTrSfwADr2_R9h-4o5mSCHccQZ1lEGNE",
    "authDomain": "filend-msg.firebaseapp.com",
    "projectId": "filend-msg",
    "storageBucket": "filend-msg.firebasestorage.app",
    "messagingSenderId": "212562196397",
    "appId": "1:212562196397:web:732a0e7dce20e6aedf9324",
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