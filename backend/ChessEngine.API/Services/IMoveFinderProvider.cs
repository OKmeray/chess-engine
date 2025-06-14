using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using ChessEngine.Application.MCTS;         // your MCTSMoveFinder
using Microsoft.Extensions.Hosting;          // IHostEnvironment
using ChessEngine.Domain.Interfaces;    // IMoveFinder

namespace ChessEngine.API.Services
{
    public interface IMoveFinderProvider : IDisposable
    {
        IMoveFinder Get(string version);
        IEnumerable<string> AvailableVersions { get; }
    }

    public class MoveFinderProvider : IMoveFinderProvider
    {
        private readonly Dictionary<string, MCTSMoveFinder> _finders;

        public MoveFinderProvider(IHostEnvironment env)
        {
            var modelsDir = Path.Combine(env.ContentRootPath, "AIModels");
            _finders = Directory
                .EnumerateFiles(modelsDir, "*.onnx")
                .ToDictionary(
                    path => Path.GetFileNameWithoutExtension(path),
                    path => new MCTSMoveFinder(path)
                );
        }

        public IMoveFinder Get(string version)
        {
            if (!_finders.TryGetValue(version, out var finder))
                throw new ArgumentException($"Unknown model version: {version}");
            return finder;
        }

        public IEnumerable<string> AvailableVersions => _finders.Keys;

        public void Dispose()
        {
            foreach (var f in _finders.Values) f.Dispose();
        }
    }
}
