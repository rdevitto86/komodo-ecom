import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/svelte';
// import Home from './home.svelte';

describe('Home Page', () => {
  // Component tests
  describe('Render', () => {
    beforeAll((suite) => {
      suite.meta.type = 'component';
    });

    it('renders home page', () => {
      // render(Home);
      
      expect(screen.getByText('Home')).toBeInTheDocument();
    });
  });

  // Integration tests
  describe('Navigation', () => {
    beforeAll((suite) => {
      suite.meta.type = 'integration';
    });

    it('navigates to Home page (self)', () => {
      // render(Home);
      
      // expect(screen.getByText('Products')).toBeInTheDocument();
    });

    it('navigates to Products page', () => {
      // render(Home);
      
      // expect(screen.getByText('Catalog')).toBeInTheDocument();
    });

    it('navigates to Services page', () => {
      // render(Home);
      
      // expect(screen.getByText('Services')).toBeInTheDocument();
    });

    it('navigates to Cart page', () => {
      // render(Home);
      
      // expect(screen.getByText('Cart')).toBeInTheDocument();
    });

    it('navigates to Profile page', () => {
      // render(Home);
      
      // expect(screen.getByText('Profile')).toBeInTheDocument();
    });

    it('navigates to Contact page', () => {
      // render(Home);
      
      // expect(screen.getByText('Contact')).toBeInTheDocument();
    });

    it('navigates to FAQ page', () => {
      // render(Home);
      
      // expect(screen.getByText('FAQ')).toBeInTheDocument();
    });
  });
});
