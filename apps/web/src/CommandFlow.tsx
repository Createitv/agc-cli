import { useState } from 'react';

/** A conceptual command flow, not a live request or execution indicator. */
export function CommandFlow({ language }: { language: 'en' | 'zh' }) {
  const [paused, setPaused] = useState(false);
  const zh = language === 'zh';
  return (
    <figure className={`commandFlow${paused ? ' isPaused' : ''}`}>
      <svg viewBox="0 0 640 310" role="img" aria-label={zh ? '终端命令连接 JSON、脚本与 AI Agent 的概念示意' : 'Conceptual terminal connections to JSON, scripts, and AI agents'}>
        <defs>
          <linearGradient id="flow-panel" x2="1" y2="1">
            <stop stopColor="#181b20" /><stop offset="1" stopColor="#0b0d10" />
          </linearGradient>
          <linearGradient id="flow-red">
            <stop stopColor="#fa454b" /><stop offset="1" stopColor="#c4313b" />
          </linearGradient>
        </defs>
        <g className="flowOrbits" fill="none" stroke="#272c33">
          <circle cx="336" cy="155" r="110" /><circle cx="336" cy="155" r="170" />
          <circle cx="336" cy="155" r="225" />
        </g>
        <g fill="none" stroke="url(#flow-red)" strokeWidth="5" strokeLinecap="round">
          <path d="M250 138 H308 C382 138 352 54 432 54 H479" />
          <path d="M250 155 H479" />
          <path d="M250 172 H308 C382 172 352 256 432 256 H479" />
        </g>
        <g className="flowPackets" fill="none" strokeWidth="6" strokeLinecap="round">
          <path className="flowPacket" pathLength="100" d="M250 138 H308 C382 138 352 54 432 54 H479" />
          <path className="flowPacket flowPacketMiddle" pathLength="100" d="M250 155 H479" />
          <path className="flowPacket flowPacketBottom" pathLength="100" d="M250 172 H308 C382 172 352 256 432 256 H479" />
        </g>
        <rect x="50" y="99" width="200" height="112" rx="23" fill="url(#flow-panel)" stroke="#646970" strokeWidth="2" />
        <path d="m87 124 25 25-25 25 M126 175h28" fill="none" stroke="#fa454b" strokeWidth="9" strokeLinecap="round" strokeLinejoin="round" />
        <text x="179" y="161" className="flowCommand">agc</text>
        <g className="flowEndpoint" fill="url(#flow-panel)" stroke="#fa454b" strokeWidth="2">
          <rect x="479" y="14" width="80" height="80" rx="20" />
          <rect x="479" y="115" width="80" height="80" rx="20" />
          <rect x="479" y="216" width="80" height="80" rx="20" />
        </g>
        <g fill="none" stroke="#d8dde3" strokeWidth="3.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M505 34h20l12 12v29h-32z M525 34v13h12 M513 56h15 M513 65h15" />
          <path d="m504 241-13 14 13 14 m30-28 13 14-13 14 m-11-30-8 30" />
          <path d="m496 140 16 15-16 15 M512 155h12 M532 152l12-13 M532 159l12 13" stroke="#fa454b" />
          <circle cx="529" cy="155" r="8" /><circle cx="547" cy="135" r="5" stroke="#62d5dd" /><circle cx="547" cy="176" r="5" stroke="#6adb9f" />
        </g>
        <g className="flowLabels">
          <text x="150" y="238" textAnchor="middle">{zh ? '终端命令' : 'Terminal'}</text>
          <text x="579" y="58">JSON</text><text x="579" y="159">Agent</text>
          <text x="579" y="260">{zh ? '脚本' : 'Script'}</text>
        </g>
        <g className="flowStatus" fill="#62d5dd" stroke="#c9f9fa" strokeWidth="3">
          <circle cx="435" cy="54" r="7" /><circle cx="435" cy="155" r="7" fill="#6adb9f" />
          <circle cx="435" cy="256" r="7" fill="#6adb9f" />
        </g>
      </svg>
      <figcaption>
        <span>{zh ? '一个命令入口，连接你的工作流。' : 'One command surface. Your workflow.'}</span>
        <button type="button" aria-pressed={paused} onClick={() => setPaused(!paused)}>
          {paused ? (zh ? '播放动效' : 'Play motion') : (zh ? '暂停动效' : 'Pause motion')}
        </button>
      </figcaption>
    </figure>
  );
}
