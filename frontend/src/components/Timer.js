import React, { useState, useEffect } from 'react';
import { Clock, Play, Pause, Square, Plus, X } from 'lucide-react';
import toast from 'react-hot-toast';

const Timer = () => {
  const [timer, setTimer] = useState({
    isRunning: false,
    duration: 0, // in seconds
    remaining: 0, // in seconds
    startTime: null
  });
  const [showTimerModal, setShowTimerModal] = useState(false);
  const [timerMinutes, setTimerMinutes] = useState(5);

  // Timer functions
  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const startTimer = () => {
    console.log('startTimer called, remaining:', timer.remaining, 'duration:', timer.duration);
    if (timer.remaining <= 0 && timer.duration <= 0) return;
    
    setTimer(prev => ({
      ...prev,
      isRunning: true,
      startTime: Date.now(),
      remaining: prev.remaining <= 0 ? prev.duration : prev.remaining
    }));
  };

  const pauseTimer = () => {
    setTimer(prev => ({
      ...prev,
      isRunning: false,
      startTime: null
    }));
  };

  const resetTimer = () => {
    setTimer({
      isRunning: false,
      duration: 0,
      remaining: 0,
      startTime: null
    });
  };

  const createTimer = () => {
    const duration = timerMinutes * 60;
    console.log('createTimer called, duration:', duration, 'minutes:', timerMinutes);
    setTimer({
      isRunning: true,
      duration: duration,
      remaining: duration,
      startTime: Date.now()
    });
    setShowTimerModal(false);
  };

  // Timer effect
  useEffect(() => {
    console.log('Timer effect running, isRunning:', timer.isRunning, 'remaining:', timer.remaining);
    let interval;
    if (timer.isRunning && timer.remaining > 0) {
      console.log('Starting interval');
      interval = setInterval(() => {
        setTimer(prev => {
          const elapsed = Math.floor((Date.now() - prev.startTime) / 1000);
          const newRemaining = Math.max(0, prev.duration - elapsed);
          
          console.log('Timer tick, elapsed:', elapsed, 'newRemaining:', newRemaining);
          
          if (newRemaining <= 0) {
            // Timer finished
            toast.success('Cronômetro finalizado!');
            return {
              ...prev,
              isRunning: false,
              remaining: 0,
              startTime: null
            };
          }
          
          return {
            ...prev,
            remaining: newRemaining
          };
        });
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [timer.isRunning, timer.startTime, timer.duration]);

  console.log('Timer component render, timer.duration:', timer.duration, 'showTimerModal:', showTimerModal);

  return (
    <>
      {/* Timer Display */}
      {timer.duration > 0 && (
        <div className="flex items-center space-x-2">
          <button
            onClick={timer.isRunning ? pauseTimer : startTimer}
            className="flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium transition-colors bg-gray-100 text-gray-700 hover:bg-gray-200"
          >
            {timer.isRunning ? (
              <Pause className="h-4 w-4" />
            ) : (
              <Play className="h-4 w-4" />
            )}
            <span>{formatTime(timer.remaining)}</span>
            {timer.isRunning && (
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            )}
          </button>
        </div>
      )}

      {/* Add Timer Button */}
      {timer.duration === 0 && (
        <button
          onClick={() => {
            console.log('Timer button clicked, showTimerModal:', showTimerModal);
            setShowTimerModal(true);
          }}
          className="flex items-center justify-center space-x-2 px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors w-32"
        >
          <span>Timer</span>
        </button>
      )}

      {/* Timer Modal */}
      {showTimerModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-[9999]">
          <div className="relative top-20 mx-auto p-6 w-11/12 md:w-1/3 bg-white rounded-lg shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-medium text-gray-900">Configurar Cronômetro</h3>
              <button
                onClick={() => setShowTimerModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Minutos
                </label>
                <div className="relative">
                  <input
                    type="number"
                    min="1"
                    max="60"
                    value={timerMinutes}
                    onChange={(e) => setTimerMinutes(Math.max(1, Math.min(60, parseInt(e.target.value) || 1)))}
                    className="w-full px-3 py-2 pr-8 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <div className="absolute inset-y-0 right-0 flex flex-col">
                    <button
                      onClick={() => setTimerMinutes(Math.min(60, timerMinutes + 1))}
                      className="flex-1 px-2 text-gray-400 hover:text-gray-600 border-l border-gray-300 rounded-r-md hover:bg-gray-50"
                    >
                      ▲
                    </button>
                    <button
                      onClick={() => setTimerMinutes(Math.max(1, timerMinutes - 1))}
                      className="flex-1 px-2 text-gray-400 hover:text-gray-600 border-l border-gray-300 rounded-r-md hover:bg-gray-50"
                    >
                      ▼
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="flex justify-end space-x-3 mt-6">
              <button
                onClick={() => setShowTimerModal(false)}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
              >
                Cancelar
              </button>
              <button
                onClick={createTimer}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Iniciar
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default Timer;
