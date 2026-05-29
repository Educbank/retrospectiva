import React, { useState, useEffect, useRef } from 'react';
import { Play, Pause, RotateCcw, ChevronDown } from 'lucide-react';
import toast from 'react-hot-toast';

const emptyTimer = {
  isRunning: false,
  duration: 0, // in seconds
  remaining: 0, // in seconds
  endTime: null // absolute ms timestamp when the timer hits 0 (only while running)
};

const PRESETS = [1, 2, 5, 10]; // minutes

// Circular progress ring sizing
const RADIUS = 18;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

const Timer = ({ retrospectiveId }) => {
  const storageKey = `educ-retro-timer-${retrospectiveId || 'default'}`;

  // Lazy init: rehydrate from localStorage, recomputing remaining for the time
  // that elapsed while the page was closed/refreshed.
  const [timer, setTimer] = useState(() => {
    try {
      const saved = JSON.parse(localStorage.getItem(storageKey));
      if (!saved || !saved.duration) return emptyTimer;

      if (saved.isRunning && saved.endTime) {
        const remaining = Math.max(0, Math.ceil((saved.endTime - Date.now()) / 1000));
        if (remaining <= 0) {
          return { ...saved, isRunning: false, remaining: 0, endTime: null };
        }
        return { ...saved, remaining };
      }
      return { ...saved, isRunning: false, endTime: null };
    } catch {
      return emptyTimer;
    }
  });
  const [showOptions, setShowOptions] = useState(false);
  const [timerMinutes, setTimerMinutes] = useState(5);
  const containerRef = useRef(null);

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const startTimer = () => {
    // When finished (remaining 0) start over from the full duration.
    const base = timer.remaining > 0 ? timer.remaining : timer.duration;
    if (base <= 0) return;

    setTimer(prev => ({
      ...prev,
      isRunning: true,
      remaining: base,
      endTime: Date.now() + base * 1000
    }));
  };

  const pauseTimer = () => {
    setTimer(prev => ({
      ...prev,
      isRunning: false,
      endTime: null
    }));
  };

  const resetTimer = () => {
    setTimer(emptyTimer);
  };

  const createTimer = (minutes) => {
    const duration = minutes * 60;
    setTimer({
      isRunning: true,
      duration: duration,
      remaining: duration,
      endTime: Date.now() + duration * 1000
    });
    setShowOptions(false);
  };

  // Persist on every change so a refresh can pick up where we left off.
  useEffect(() => {
    try {
      if (timer.duration > 0) {
        localStorage.setItem(storageKey, JSON.stringify(timer));
      } else {
        localStorage.removeItem(storageKey);
      }
    } catch {
      // ignore storage errors (e.g. private mode)
    }
  }, [timer, storageKey]);

  // Countdown tick driven by the absolute endTime.
  useEffect(() => {
    if (!timer.isRunning || !timer.endTime) return;

    const interval = setInterval(() => {
      setTimer(prev => {
        const remaining = Math.max(0, Math.ceil((prev.endTime - Date.now()) / 1000));
        if (remaining <= 0) {
          toast.success('Cronômetro finalizado!');
          return { ...prev, isRunning: false, remaining: 0, endTime: null };
        }
        return { ...prev, remaining };
      });
    }, 250);

    return () => clearInterval(interval);
  }, [timer.isRunning, timer.endTime]);

  // Close the dropdown when clicking outside.
  useEffect(() => {
    if (!showOptions) return;

    const handleClickOutside = (e) => {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setShowOptions(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [showOptions]);

  // Ring state
  const progress = timer.duration > 0 ? timer.remaining / timer.duration : 0;
  const finished = timer.duration > 0 && timer.remaining === 0;
  const low = timer.isRunning && timer.remaining > 0 && timer.remaining <= 10;
  const ringColor = finished ? '#ef4444' : low ? '#f59e0b' : '#2563eb';
  const textColor = finished ? 'text-red-600' : low ? 'text-amber-600' : 'text-gray-800';

  return (
    <div className="relative" ref={containerRef}>
      {/* Active timer - circular ring (EasyRetro style) */}
      {timer.duration > 0 && (
        <div className="flex items-center space-x-2">
          <div className={`relative h-11 w-11 ${finished ? 'animate-pulse' : ''}`}>
            <svg className="h-11 w-11 -rotate-90" viewBox="0 0 44 44">
              <circle cx="22" cy="22" r={RADIUS} fill="none" stroke="#e5e7eb" strokeWidth="4" />
              <circle
                cx="22"
                cy="22"
                r={RADIUS}
                fill="none"
                stroke={ringColor}
                strokeWidth="4"
                strokeLinecap="round"
                strokeDasharray={CIRCUMFERENCE}
                strokeDashoffset={CIRCUMFERENCE * (1 - progress)}
                style={{ transition: 'stroke-dashoffset 0.3s linear, stroke 0.3s ease' }}
              />
            </svg>
            <span
              className={`absolute inset-0 flex items-center justify-center text-[11px] font-semibold tabular-nums ${textColor}`}
            >
              {formatTime(timer.remaining)}
            </span>
          </div>

          <button
            onClick={timer.isRunning ? pauseTimer : startTimer}
            className="flex items-center justify-center p-2 rounded-md text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
            title={timer.isRunning ? 'Pausar' : timer.remaining > 0 ? 'Continuar' : 'Reiniciar'}
          >
            {timer.isRunning ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
          </button>
          <button
            onClick={resetTimer}
            className="flex items-center justify-center p-2 rounded-md text-gray-700 bg-gray-100 hover:bg-red-100 hover:text-red-600 transition-colors"
            title="Parar e remover cronômetro"
          >
            <RotateCcw className="h-4 w-4" />
          </button>
        </div>
      )}

      {/* Combo trigger - opens the options dropdown below */}
      {timer.duration === 0 && (
        <button
          onClick={() => setShowOptions((v) => !v)}
          className="flex items-center justify-between space-x-2 px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors w-32"
        >
          <span>Timer</span>
          <ChevronDown className={`h-4 w-4 transition-transform ${showOptions ? 'rotate-180' : ''}`} />
        </button>
      )}

      {/* Options dropdown */}
      {showOptions && timer.duration === 0 && (
        <div className="absolute right-0 mt-2 w-56 bg-white border border-gray-200 rounded-lg shadow-lg z-50 p-3">
          <p className="text-xs font-medium text-gray-500 mb-2">Iniciar cronômetro</p>

          {/* Quick presets */}
          <div className="grid grid-cols-2 gap-2 mb-3">
            {PRESETS.map((min) => (
              <button
                key={min}
                onClick={() => createTimer(min)}
                className="px-3 py-2 text-sm font-medium rounded-md border border-gray-200 text-gray-700 hover:bg-blue-50 hover:border-blue-300 hover:text-blue-700 transition-colors"
              >
                {min} min
              </button>
            ))}
          </div>

          {/* Custom minutes */}
          <div className="border-t border-gray-100 pt-3">
            <label className="block text-xs font-medium text-gray-500 mb-1">Personalizado</label>
            <div className="flex items-center space-x-2">
              <input
                type="number"
                min="1"
                max="60"
                value={timerMinutes}
                onChange={(e) => setTimerMinutes(Math.max(1, Math.min(60, parseInt(e.target.value) || 1)))}
                className="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
              <button
                onClick={() => createTimer(timerMinutes)}
                className="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors whitespace-nowrap"
              >
                Iniciar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Timer;
