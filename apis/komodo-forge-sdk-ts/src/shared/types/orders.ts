export interface Order {
  id: string;
  userId: string;
  status: 'pending' | 'confirmed' | 'shipped' | 'delivered' | 'cancelled';
  total: number;
  currency: 'USD';
  items: OrderItem[];
  createdAt: string;
  updatedAt: string;
}

export interface OrderItem {
  serviceId: string;
  quantity: number;
  price: number;
  currency: 'USD';
  metadata?: Record<string, any>;
}