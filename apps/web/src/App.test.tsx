import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { App } from './App';

describe('App', () => {
  it('renders the release workflow, profile binding, and API registry', () => {
    render(<App />);
    expect(screen.getByText('agccli.app')).toBeInTheDocument();
    expect(screen.getByText(/Move every release through/i)).toBeInTheDocument();
    expect(screen.getByText('Automatic profile binding')).toBeInTheDocument();
    expect(screen.getAllByText('Publishing API').length).toBeGreaterThan(0);
    expect(screen.getByText(/registered interfaces/)).toBeInTheDocument();
  });
});
