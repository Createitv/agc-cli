import '@testing-library/jest-dom/vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { App } from './App';

describe('App', () => {
  afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

  beforeEach(() => {
    window.history.replaceState({}, '', '/en/');
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

    expect(screen.getByRole('dialog', { name: '像普通 CLI 工具一样安装 agc。' })).toBeInTheDocument();
    expect(screen.getAllByText(/brew install agc-cli/).length).toBeGreaterThan(1);
    expect(screen.getAllByText(/scoop install agc-cli/).length).toBeGreaterThan(0);
    expect(document.documentElement).toHaveAttribute('lang', 'zh-CN');
  });

  it('shows the complete selected family without a local API server', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')));
    render(<App />);
    expect(await screen.findByText('app-info-query')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /Reports 12 interfaces/ }));
    expect(screen.queryByText('app-info-query')).not.toBeInTheDocument();
    expect(document.querySelectorAll('.endpointEntry')).toHaveLength(12);
    fireEvent.click(screen.getByRole('button', { name: /Testing 27 interfaces/ }));
    expect(document.querySelectorAll('.endpointEntry')).toHaveLength(27);
  });

  it('uses the URL locale for Chinese pages even when English was stored', () => {
    window.history.replaceState({}, '', '/reference/reports/');
    window.localStorage.setItem('agc-language', 'en');
    render(<App />);
    expect(document.documentElement).toHaveAttribute('lang', 'zh-CN');
    expect(document.querySelectorAll('.endpointEntry')).toHaveLength(12);
    fireEvent.click(screen.getByRole('button', { name: '语言' }));
    fireEvent.click(screen.getByRole('menuitemradio', { name: /English/ }));
    expect(window.location.pathname).toBe('/en/reference/reports/');
    expect(document.documentElement).toHaveAttribute('lang', 'en');
    expect(document.querySelectorAll('.endpointEntry')).toHaveLength(12);
  });
});
