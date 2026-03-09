import { describe, it, expect, beforeAll } from 'vitest';
import { render, screen } from '@testing-library/svelte';
// import Error from './+error.svelte';


describe('Error Page', () => {
  beforeAll((suite) => {
    suite.meta.type = 'component';
  });

  it('renders error message', () => {
    // render(Error, { status: 404, error: new Error('Not Found'), message: 'Not Found' });
    
    expect(screen.getByText('Not Found')).toBeInTheDocument();
  });
});
