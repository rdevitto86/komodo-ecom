export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  duration?: number;
}

class NotificationsState {
  #notifications = $state<Notification[]>([]);
  #toastPending = $state<Notification[]>([]);
  #snackbarPending = $state<Notification[]>([]);
  
  get notificationsCount() { return this.#notifications.length; }
  get toastNotificationsCount() { return this.#toastPending.length; }
  get snackbarNotificationsCount() { return this.#snackbarPending.length; }

  addNotification = (notification: Notification) => {
    // TODO
    this.#notifications.push(notification);
  };
  addToastNotification = (notification: Notification) => {
    // TODO
    this.#toastPending.push(notification);
  };
  addSnackbarNotification = (notification: Notification) => {
    // TODO
    this.#snackbarPending.push(notification);
  };
}

export const notificationsState = new NotificationsState();
