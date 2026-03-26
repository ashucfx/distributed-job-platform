"use client";

import { useEffect, useState, useCallback } from "react";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function Home() {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [jobName, setJobName] = useState("");
  const [payload, setPayload] = useState("{\n  \"to\": \"user@example.com\",\n  \"subject\": \"Hello from Platform\"\n}");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchJobs = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/api/v1/jobs`);
      if (res.ok) {
        const data = await res.json();
        setJobs(data || []);
      }
    } catch (e) {
      console.error("Failed to fetch jobs:", e);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchJobs();
    const interval = setInterval(fetchJobs, 5000);
    return () => clearInterval(interval);
  }, [fetchJobs]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      let parsedPayload;
      try {
        parsedPayload = JSON.parse(payload);
      } catch (e) {
        alert("Invalid JSON payload");
        setIsSubmitting(false);
        return;
      }

      const res = await fetch(`${API_BASE_URL}/api/v1/jobs`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: jobName,
          payload: parsedPayload,
        }),
      });

      if (res.ok) {
        setJobName("");
        setPayload("{\n  \"to\": \"user@example.com\",\n  \"subject\": \"Hello from Platform\"\n}");
        fetchJobs();
      } else {
        const errData = await res.json();
        alert(`Error: ${errData.error || "Failed to submit job"}`);
      }
    } catch (e) {
      console.error("Submission error:", e);
      alert("Network error. Is the API running?");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen p-4 md:p-8">
      <div className="max-w-7xl mx-auto space-y-10">
        {/* Header */}
        <header className="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="space-y-2">
            <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-emerald-500">
              Distributed Job Platform
            </h1>
            <p className="text-gray-500 dark:text-gray-400 font-medium text-lg">
              Enterprise-grade job orchestration & monitoring
            </p>
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-sm font-semibold animate-pulse-slow">
              <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
              System Live
            </div>
            <a 
              href="http://localhost:3002" 
              target="_blank" 
              rel="noopener noreferrer"
              className="px-5 py-2.5 rounded-xl bg-gray-900 dark:bg-white dark:text-gray-900 text-white font-bold hover:scale-105 transition-all shadow-xl shadow-blue-500/10"
            >
              Analytics Dashboard
            </a>
          </div>
        </header>

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          {/* Submission Form */}
          <section className="lg:col-span-4 glass rounded-3xl p-8 shadow-2xl border-white/20">
            <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
              <svg className="w-6 h-6 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              Quick Enqueue
            </h2>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-bold uppercase tracking-wider text-gray-400">Job Type</label>
                <input 
                  type="text" 
                  value={jobName} 
                  onChange={(e) => setJobName(e.target.value)} 
                  className="w-full bg-gray-100/50 dark:bg-gray-800/50 border-0 rounded-2xl p-4 focus:ring-2 focus:ring-blue-500 transition-all font-medium"
                  placeholder="e.g. send_email"
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-bold uppercase tracking-wider text-gray-400">Payload Parameters</label>
                <textarea 
                  value={payload} 
                  onChange={(e) => setPayload(e.target.value)} 
                  className="w-full bg-gray-100/50 dark:bg-gray-800/50 border-0 rounded-2xl p-4 font-mono text-sm h-48 focus:ring-2 focus:ring-blue-500 transition-all"
                  required
                />
              </div>
              <button 
                type="submit" 
                disabled={isSubmitting}
                className={`w-full py-4 rounded-2xl font-black text-lg transition-all shadow-lg 
                  ${isSubmitting ? 'bg-gray-400 cursor-not-allowed' : 'bg-gradient-to-r from-blue-600 to-blue-500 text-white hover:shadow-blue-500/40 hover:-translate-y-1'}`}
              >
                {isSubmitting ? "Dispatching..." : "Dispatch Job"}
              </button>
            </form>
          </section>

          {/* Job Table */}
          <section className="lg:col-span-8 glass rounded-3xl overflow-hidden shadow-2xl border-white/20">
            <div className="p-8 border-b border-white/10 flex items-center justify-between">
              <h2 className="text-2xl font-bold flex items-center gap-2">
                <svg className="w-6 h-6 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                Execution Monitor
              </h2>
              <button onClick={() => fetchJobs()} className="text-blue-500 hover:text-blue-400 font-bold transition-colors">
                Refresh Now
              </button>
            </div>
            
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-gray-500/5 text-xs font-black uppercase tracking-widest text-gray-400">
                    <th className="p-6">Job Identifier</th>
                    <th className="p-6">Execution Status</th>
                    <th className="p-6 text-center">Attempts</th>
                    <th className="p-6">Timestamp</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5">
                  {loading && jobs.length === 0 ? (
                    <tr><td colSpan={4} className="p-20 text-center text-gray-500 font-bold animate-pulse">Synchronizing with Cluster...</td></tr>
                  ) : jobs.length === 0 ? (
                    <tr><td colSpan={4} className="p-20 text-center text-gray-500 font-bold">Waiting for workload... Enqueue a job to begin tracking.</td></tr>
                  ) : (
                    jobs.map((job: any) => (
                      <tr key={job.id} className="hover:bg-blue-500/5 transition-colors group">
                        <td className="p-6">
                          <div className="flex flex-col">
                            <span className="font-black text-gray-900 dark:text-white group-hover:text-blue-500 transition-colors uppercase text-sm">{job.name}</span>
                            <span className="font-mono text-[10px] text-gray-500 truncate w-32">{job.id}</span>
                          </div>
                        </td>
                        <td className="p-6">
                          <StatusBadge status={job.status} />
                        </td>
                        <td className="p-6">
                          <div className="flex flex-col items-center gap-1">
                            <div className="flex gap-1">
                              {[...Array(job.max_retries)].map((_, i) => (
                                <div key={i} className={`w-2 h-2 rounded-full ${i < job.retries ? 'bg-red-400' : 'bg-gray-200 dark:bg-gray-700'}`}></div>
                              ))}
                            </div>
                            <span className="text-[10px] font-bold text-gray-400">{job.retries}/{job.max_retries}</span>
                          </div>
                        </td>
                        <td className="p-6">
                          <span className="font-medium text-gray-400 text-sm whitespace-nowrap">
                            {new Date(job.created_at).toLocaleString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
                          </span>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  const configs: Record<string, { bg: string, text: string, label: string }> = {
    'COMPLETED': { bg: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-400', text: 'text-emerald-500', label: 'Success' },
    'FAILED': { bg: 'bg-rose-100 text-rose-700 dark:bg-rose-500/20 dark:text-rose-400', text: 'text-rose-500', label: 'Failure' },
    'PROCESSING': { bg: 'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-400', text: 'text-blue-500', label: 'Executing' },
    'PENDING': { bg: 'bg-gray-100 text-gray-700 dark:bg-gray-700/50 dark:text-gray-400', text: 'text-gray-400', label: 'Queued' },
  };

  const config = configs[status] || configs['PENDING'];

  return (
    <div className={`inline-flex items-center gap-2 px-4 py-1.5 rounded-xl font-black text-[11px] uppercase tracking-tighter shadow-sm ${config.bg}`}>
      {status === 'PROCESSING' && <span className="w-1.5 h-1.5 rounded-full bg-blue-500 animate-ping"></span>}
      {config.label}
    </div>
  );
}
