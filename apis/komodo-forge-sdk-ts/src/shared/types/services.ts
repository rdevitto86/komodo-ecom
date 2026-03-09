export interface Service {
  id: string;
  slug: string;
  sku?: string;
  name: string;
  description: string;
  category: ServiceCategory;
  status: 'active' | 'draft' | 'archived';
  price: number;
  compareAtPrice?: number;
  cost?: number;
  currency?: 'USD';
  taxCode?: string;
  duration?: ServiceDuration;
  images?: ServiceImage[];
  serviceOptions?: ServiceOption[];
  locationTypes: LocationType[];
  availability?: ServiceAvailability;
  requirements?: string[];
  includedItems?: string[];
  relatedServiceIds?: string[];
  relatedProductIds?: string[];
  meta?: {
    tags?: string[];
    isPopular?: boolean;
    isFeatured?: boolean;
    isNew?: boolean;
    [key: string]: any;
  };
  seo?: {
    title?: string;
    description?: string;
    keywords?: string[];
  };
  createdAt?: string;
  updatedAt?: string;
}

export interface ServiceImage {
  url: string;
  alt?: string;
  isPrimary?: boolean;
  type?: 'image' | 'video';
}

export type ServiceCategory = 
  | 'installation' 
  | 'repair' 
  | 'maintenance' 
  | 'consultation' 
  | 'delivery' 
  | 'setup'
  | 'custom';

export type LocationType = 'residential' | 'commercial' | 'industrial' | 'remote';

export interface ServiceDuration {
  estimated: number;
  unit: 'minutes' | 'hours' | 'days';
  min?: number;
  max?: number;
}

export interface ServiceOption {
  id: string;
  name: string;
  description?: string;
  type: 'addon' | 'upgrade' | 'required';
  priceModifier: number;
  durationModifier?: number;
  maxBookings?: number;
  isDefault?: boolean;
  requiresProducts?: string[];
  compatibleWith?: string[];
}

export interface ServiceAvailability {
  daysOfWeek?: number[];
  timeSlots?: TimeSlot[];
  blackoutDates?: string[];
  leadTimeDays?: number;
  maxBookingsPerDay?: number;
  serviceAreaZipCodes?: string[];
  serviceAreaRadius?: {
    miles: number;
    centerZip: string;
  };
}

export interface TimeSlot {
  start: string;
  end: string;
  available: boolean;
  maxBookings?: number;
}