import '@testing-library/jest-dom/vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { App } from './App';

describe('App', () => {
  afterEach(() => cleanup());

  beforeEach(() => {
    const values = new Map<string, string>();
    Object.defineProperty(window, 'localStorage', {
      configurable: true,
      value: {
        getItem: (key: string) => values.get(key) ?? null,
        setItem: (key: string, value: string) => values.set(key, value),
        removeItem: (key: string) => values.delete(key),
        clear: () => values.clear(),
      },
    });
  });

  it('renders the release workflow, profile binding, and API registry', () => {
    render(<App />);
    expect(screen.getByText('agccli.app')).toBeInTheDocument();
    expect(screen.getByText(/Move every release through/i)).toBeInTheDocument();
    expect(screen.getByText('Automatic profile binding')).toBeInTheDocument();
    expect(screen.getAllByText('Publishing API').length).toBeGreaterThan(0);
    expect(screen.getByText(/registered interfaces/)).toBeInTheDocument();
  });

  it('switches to Chinese and opens the direct install dialog', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Language' }));
    fireEvent.click(screen.getByRole('menuitemradio', { name: /简体中文/ }));

    expect(screen.getByText(/让每一次发布沿着/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: '安装' }));

    expect(screen.getByRole('dialog', { name: '把 agc 放进 shell 的 PATH。' })).toBeInTheDocument();
    expect(screen.getByText(/go install github\.com\/Createitv\/agc-cli\/cmd\/agc@latest/)).toBeInTheDocument();
    expect(document.documentElement).toHaveAttribute('lang', 'zh-CN');
  });
});
