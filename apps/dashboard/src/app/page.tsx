"use client";

import { useEffect, useState } from "react";

export default function Home() {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [jobName, setJobName] = useState("");
  const [payload, setPayload] = useState("{}");

  useEffect(() => {
    fetchJobs();
    const interval = setInterval(fetchJobs, 3000);
    return () => clearInterval(interval);
  }, []);

  const fetchJobs = async () => {
    try {
      // In a real app we'd have a List endpoint, for now we simulate or rely on standard queries
      // We will add a simple mock or if you extended the API to return a list.
      // Assuming GET /api/v1/jobs exists and returns an array
      const res = await fetch("http://localhost:8080/api/v1/jobs");
      if (res.ok) {
        const data = await res.json();
        setJobs(data);
      }
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await fetch("http://localhost:8080/api/v1/jobs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: jobName,
          payload: JSON.parse(payload),
        }),
      });
      setJobName("");
      setPayload("{}");
      fetchJobs();
    } catch (e) {
      alert("Failed to submit job. Check console.");
      console.error(e);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 p-8 font-sans">
      <div className="max-w-6xl mx-auto space-y-8">
        <header className="flex items-center justify-between border-b pb-6">
          <div>
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">Distributed Job Platform</h1>
            <p className="text-gray-500 mt-2">Real-time job queue monitoring and submission</p>
          </div>
          <div className="flex gap-4">
            <a href="http://localhost:3002/d/job-metrics" target="_blank" className="bg-orange-500 text-white px-4 py-2 rounded shadow hover:bg-orange-600 transition">View Grafana</a>
          </div>
        </header>

        <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white p-6 rounded-lg shadow border col-span-1">
            <h2 className="text-xl font-semibold mb-4">Submit New Job</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Job Name</label>
                <input 
                  type="text" 
                  value={jobName} 
                  onChange={(e) => setJobName(e.target.value)} 
                  className="w-full border rounded p-2 text-sm"
                  placeholder="e.g. send_email"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Payload (JSON)</label>
                <textarea 
                  value={payload} 
                  onChange={(e) => setPayload(e.target.value)} 
                  className="w-full border rounded p-2 font-mono text-sm h-32"
                  required
                />
              </div>
              <button type="submit" className="w-full bg-black text-white px-4 py-2 rounded hover:bg-gray-800 transition">
                Enqueue Job
              </button>
            </form>
          </div>

          <div className="bg-white p-6 rounded-lg shadow border col-span-1 md:col-span-2">
            <h2 className="text-xl font-semibold mb-4">Recent Jobs</h2>
            {loading ? (
              <div className="text-center text-gray-500 py-10">Loading jobs...</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="p-3 font-medium">ID</th>
                      <th className="p-3 font-medium">Name</th>
                      <th className="p-3 font-medium">Status</th>
                      <th className="p-3 font-medium">Retries</th>
                      <th className="p-3 font-medium">Created At</th>
                    </tr>
                  </thead>
                  <tbody>
                    {jobs.length === 0 ? (
                      <tr><td colSpan={5} className="p-3 text-center text-gray-500">No jobs found in the database. Wait for scheduler or submit one.</td></tr>
                    ) : (
                      jobs.map((job: any) => (
                        <tr key={job.id} className="border-b hover:bg-gray-50">
                          <td className="p-3 font-mono text-xs">{job.id.substring(0, 8)}...</td>
                          <td className="p-3">{job.name}</td>
                          <td className="p-3">
                            <span className={`px-2 py-1 rounded-full text-xs font-medium 
                              ${job.status === 'COMPLETED' ? 'bg-green-100 text-green-800' : 
                                job.status === 'FAILED' ? 'bg-red-100 text-red-800' : 
                                job.status === 'PROCESSING' ? 'bg-blue-100 text-blue-800' : 
                                'bg-gray-100 text-gray-800'}`}>
                              {job.status}
                            </span>
                          </td>
                          <td className="p-3">{job.retries} / {job.max_retries}</td>
                          <td className="p-3">{new Date(job.created_at).toLocaleTimeString()}</td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </section>
      </div>
    </div>
  );
}
