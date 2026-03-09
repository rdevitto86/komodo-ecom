export interface Address {
  line1: string;
  line2?: string;
  city: string;
  state: string;
  zip: string;
  country: string;
}

export interface UserPreferences {
  theme: 'light' | 'dark';
  language: string;
  currency: string;
}

export interface Orders {

}

export interface UserProfile {
  id: string;
  name: string;
  email: string;
  phone?: string;
  avatar?: string;
  role: 'GUEST' | 'CUSTOMER';
  address?: Address[];
  payments?: string[];
  preferences?: UserPreferences;
  orders?: Orders[];
  createdAt: Date;
  meta?: Record<string, any>;
}

class UserState {
  #profile = $state<UserProfile | null>(null);
  #isAuthorized = $state(false);
  
  get profile() { return this.#profile; }
  get isAuthorized() { return this.#isAuthorized; }

  login() {
    this.#isAuthorized = true;
  }

  logout() {
    this.#isAuthorized = false;
    this.#profile = null;
  }
}

export const userState = new UserState();