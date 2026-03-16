import React, { useState, useEffect } from "react";
import { Clock, Car, User, Smartphone, Power } from "lucide-react";
import {
  formatDistanceToNow,
  isAfter,
  addMinutes,
  differenceInSeconds,
} from "date-fns";

const BookingDetails = ({ booking, spot, onEndSession }) => {
  const [timeLeft, setTimeLeft] = useState("");
  const [isOvertime, setIsOvertime] = useState(false);
  const [progress, setProgress] = useState(100);

  const endTime = addMinutes(new Date(booking.startTime), booking.duration);

  useEffect(() => {
    const timer = setInterval(() => {
      const now = new Date();
      const expired = isAfter(now, endTime);
      setIsOvertime(expired);

      const totalSeconds = booking.duration * 60;
      const remainingSeconds = Math.max(0, differenceInSeconds(endTime, now));
      const percentage = (remainingSeconds / totalSeconds) * 100;
      setProgress(expired ? 0 : percentage);

      if (expired) {
        setTimeLeft(formatDistanceToNow(endTime) + " (Overtime)");
      } else {
        setTimeLeft(
          formatDistanceToNow(now, { addSuffix: false }) + " tersisa",
        );
      }
    }, 1000);

    return () => clearInterval(timer);
  }, [booking, endTime]);

  return (
    <div
      className={`p-6 lg:p-8 rounded-3xl border ${isOvertime ? "bg-red-50/50 border-red-100" : "bg-white border-slate-200"} shadow-sm space-y-8`}
    >
      <div className="flex justify-between items-start">
        <div className="flex items-center gap-3">
          <div
            className={`p-3 rounded-2xl ${isOvertime ? "bg-red-100 text-red-600" : "bg-primary/10 text-primary"}`}
          >
            <Car size={24} />
          </div>
          <div>
            <h4 className="text-2xl font-black font-heading tracking-tight">
              {spot.label}
            </h4>
            <span
              className={`text-[10px] font-black uppercase px-2 py-0.5 rounded-md ${isOvertime ? "bg-red-500 text-white" : "bg-success text-white"}`}
            >
              {isOvertime ? "Overtime" : "Proses"}
            </span>
          </div>
        </div>
        <button
          onClick={() => onEndSession(booking.id)}
          className="p-3 text-slate-400 hover:text-danger hover:bg-danger/10 rounded-xl transition-all"
          title="Akhiri Sesi"
        >
          <Power size={20} />
        </button>
      </div>

      <div className="space-y-4">
        <div className="flex justify-between items-end">
          <div className="flex items-center gap-2 text-slate-500">
            <Clock size={16} />
            <span className="text-sm font-semibold italic">{timeLeft}</span>
          </div>
          <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">
            {Math.round(progress)}%
          </span>
        </div>
        <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
          <div
            className={`h-full transition-all duration-1000 ease-linear rounded-full ${isOvertime ? "bg-danger" : "bg-primary"}`}
            style={{ width: `${progress}%` }}
          ></div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6 bg-slate-50 p-5 rounded-2xl">
        <div className="space-y-1">
          <span className="text-[10px] font-bold text-slate-400 uppercase block tracking-tighter">
            Pengendara
          </span>
          <div className="flex items-center gap-2 text-slate-900">
            <User size={14} className="text-slate-400" />
            <span className="font-bold text-sm truncate">{booking.name}</span>
          </div>
        </div>
        <div className="space-y-1">
          <span className="text-[10px] font-bold text-slate-400 uppercase block tracking-tighter">
            No. Kendaraan
          </span>
          <div className="flex items-center gap-2 text-slate-900">
            <Smartphone size={14} className="text-slate-400" />
            <span className="font-bold text-sm truncate">
              {booking.vehicleNumber}
            </span>
          </div>
        </div>
      </div>

      <button
        className="w-full py-4 bg-slate-900 text-white rounded-2xl font-bold hover:bg-slate-800 transition-all flex items-center justify-center gap-3 shadow-xl shadow-slate-200"
        onClick={() => onEndSession(booking.id)}
      >
        <CheckCircle size={20} /> Checkout Sesi
      </button>
    </div>
  );
};

// Simple internal icon if CheckCircle is missing from imports
const CheckCircle = ({ size }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="3"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
    <polyline points="22 4 12 14.01 9 11.01"></polyline>
  </svg>
);

export default BookingDetails;
