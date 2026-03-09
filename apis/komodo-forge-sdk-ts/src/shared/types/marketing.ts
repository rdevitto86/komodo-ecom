export type MarketingContentType = 
  | 'hero-banner' | 'inline-banner' | 'promo-card' | 'recommendation-grid'
  | 'countdown-offer' | 'social-proof' | 'cross-sell' | 'announcement-bar'
  | 'modal-popup' | 'exit-intent' | 'sticky-footer';

export type ContentFormat = 'structured' | 'html' | 'markdown';

export interface MarketingContent {
  id: string;
  code: string;
  name: string;
  type: MarketingContentType;
  status: 'draft' | 'scheduled' | 'active' | 'paused' | 'expired';
  format: ContentFormat;
  priority: number;
  placement?: 'top' | 'middle' | 'bottom' | 'sidebar' | 'overlay' | 'inline';
  campaignId?: string;
  variantId?: string;
  content: MarketingContentData;
  targeting?: MarketingTargeting;
  schedule?: {
    startDate?: string;
    endDate?: string;
    timezone?: string;
  };
  tracking?: {
    impressionEvent?: string;
    clickEvent?: string;
    conversionEvent?: string;
  };
  createdAt?: string;
  updatedAt?: string;
}

export interface MarketingContentData {
  title?: string;
  subtitle?: string;
  description?: string;
  imageUrl?: string;
  imageAlt?: string;
  backgroundImageUrl?: string;
  videoUrl?: string;
  ctaText?: string;
  ctaUrl?: string;
  secondaryCtaText?: string;
  secondaryCtaUrl?: string;
  htmlContent?: string;
  markdownContent?: string;
  badgeText?: string;
  countdownEndDate?: string;
  productIds?: string[];
  serviceIds?: string[];
  discountCode?: string;
  [key: string]: any;
}

export interface MarketingTargeting {
  userSegments?: string[];
  excludeSegments?: string[];
  locations?: string[];
  excludeLocations?: string[];
  deviceTypes?: ('mobile' | 'tablet' | 'desktop')[];
  newUsersOnly?: boolean;
  returningUsersOnly?: boolean;
  minCartValue?: number;
  maxCartValue?: number;
  hasProducts?: string[];
  viewedProducts?: string[];
  categories?: string[];
  urlPatterns?: string[];
  excludeUrlPatterns?: string[];
}

export interface Campaign {
  id: string;
  name: string;
  description?: string;
  status: 'draft' | 'active' | 'paused' | 'completed';
  startDate?: string;
  endDate?: string;
  budget?: number;
  contentIds: string[];
  abTestEnabled?: boolean;
  createdAt?: string;
  updatedAt?: string;
}