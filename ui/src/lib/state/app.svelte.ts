export interface PageSnapshot {
  url: string;
  scrollY: number;
  state: any;
  rehydrate: boolean;
}

class AppState {
  #version = import.meta.env.VERSION;
  #sidebarMin = $state(false);
  #isLoading = $state(false);
  #showAnimations = $state(true);

  #sessionId = $state<string | null>(null);
  #screenSize = $state<'mobile' | 'tablet' | 'desktop'>('desktop');
  #theme = $state<'light' | 'dark'>('light');
  #language = $state('en');
  #currency = $state('USD');

  #landingSlotSeed = $state<string | null>(null);
  #breadcrumbs = $state<PageSnapshot[]>([]);
  
  // Computed Props
  isMobile = $derived(this.#screenSize === 'mobile');
  isTablet = $derived(this.#screenSize === 'tablet');
  isDesktop = $derived(this.#screenSize === 'desktop');
  isScrollLocked = $derived(this.#isLoading || (this.#sidebarMin && this.isMobile));
  currentPath = $derived(this.#breadcrumbs[this.#breadcrumbs.length - 1]?.url ?? '/');

  constructor() {
    if (typeof window !== 'undefined') {
      Object.entries({
        mobile: '(max-width: 767px)',
        tablet: '(min-width: 768px) and (max-width: 1024px)',
        desktop: '(min-width: 1025px)'
      }).forEach(([platform, query]) => {
        const media = window.matchMedia(query);
        const handler = () => media.matches && (this.#screenSize = platform as any);
        media.addEventListener('change', handler);
        handler();
      });

      const savedTheme = localStorage.getItem('theme') as 'light' | 'dark';
      if (savedTheme) this.#theme = savedTheme;

      $effect.root(() => {
        $effect(() => {
          localStorage.setItem('theme', this.#theme);
          document.documentElement.setAttribute('data-theme', this.#theme);
          document.documentElement.classList.toggle('dark', this.#theme === 'dark');
        });

        $effect(() => {
          if (this.isScrollLocked) {
            document.body.style.overflow = 'hidden';
            document.body.style.paddingRight = `${window.innerWidth - document.documentElement.clientWidth}px`;
          } else {
            document.body.style.overflow = '';
            document.body.style.paddingRight = '';
          }
        });
      });
    }
  }

  // Getters
  get version() { return this.#version; }
  get sidebarMin() { return this.#sidebarMin; }
  get isLoading() { return this.#isLoading; }
  get showAnimations() { return this.#showAnimations; }
  get screenSize() { return this.#screenSize; }
  get sessionId() { return this.#sessionId; }
  get theme() { return this.#theme; }
  get language() { return this.#language; }
  get currency() { return this.#currency; }
  get landingSlotSeed() { return this.#landingSlotSeed; }
  get breadcrumbs() { return this.#breadcrumbs; }

  // Setters
  setIsLoading = (loading: boolean) => this.#isLoading = loading;
  setShowAnimations = (show: boolean) => this.#showAnimations = show;
  setTheme = (theme: 'light' | 'dark') => this.#theme = theme;
  setLanguage = (language: string) => this.#language = language;
  setCurrency = (currency: string) => this.#currency = currency;
  setSessionId = (sessionId: string) => this.#sessionId = sessionId;
  setLandingSlotSeed = (seed: string) => this.#landingSlotSeed = seed;

  // Actions
  toggleSidebar = () => this.#sidebarMin = !this.#sidebarMin;
  restoreBreadcrumbState = (url: string) => {
    if (url && this.#breadcrumbs.some((b) => b.url === url)) {
      // TODO: Implement breadcrumb state restoration
    }
  };

  // Navigation
  addBreadcrumb = (breadcrumb: PageSnapshot) => this.#breadcrumbs.push(breadcrumb);
  removeBreadcrumb = (url: string) => (this.#breadcrumbs = this.#breadcrumbs.filter((b) => b.url !== url));
  clearBreadcrumbs = () => (this.#breadcrumbs = []);
}

export const appState = new AppState();
