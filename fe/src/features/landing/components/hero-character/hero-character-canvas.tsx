'use client'

import { Suspense, useLayoutEffect } from 'react'
import { Canvas, useThree } from '@react-three/fiber'
import { ContactShadows, Environment } from '@react-three/drei'
import * as THREE from 'three'

import { GltfHeroCharacter } from './gltf-hero-character'

function FramedCamera() {
  const camera = useThree((s) => s.camera)
  const size = useThree((s) => s.size)

  useLayoutEffect(() => {
    // Pull back so full body fits with padding inside the frame
    camera.position.set(0, 1.0, 4.6)
    camera.near = 0.1
    camera.far = 50
    if (camera instanceof THREE.PerspectiveCamera) {
      camera.fov = 30
      camera.aspect = size.width / Math.max(size.height, 1)
      camera.updateProjectionMatrix()
    }
    camera.lookAt(0, 0.9, 0)
  }, [camera, size.height, size.width])

  return null
}

function TransparentScene() {
  const { gl, scene } = useThree()

  useLayoutEffect(() => {
    // Three.js has no "transparent" Color — that string becomes white.
    scene.background = null
    gl.setClearColor(0x000000, 0)
  }, [gl, scene])

  return null
}

function CharacterScene() {
  return (
    <>
      <TransparentScene />
      <FramedCamera />
      <ambientLight intensity={0.7} />
      <directionalLight castShadow position={[2.8, 5.5, 3]} intensity={1.05} />
      <directionalLight position={[-2.5, 1.5, -1.5]} intensity={0.25} />
      <hemisphereLight args={['#dbeafe', '#1e293b', 0.35]} />

       <GltfHeroCharacter /> 

      <ContactShadows position={[0, 0.01, 0]} opacity={0.35} scale={6} blur={2.2} far={3.5} />
      <Environment preset="apartment" environmentIntensity={0.2} />
    </>
  )
}

type HeroCharacterCanvasProps = {
  className?: string
}

/** WebGL stage for the hero character. Parent handles reduced-motion fallback. */
export function HeroCharacterCanvas({ className }: HeroCharacterCanvasProps) {
  return (
    <div className={className}>
      {/* <Canvas
        shadows
        dpr={[1, 1.75]}
        camera={{ position: [0, 1.0, 4.6], fov: 30 }}
        gl={{ antialias: true, alpha: true, premultipliedAlpha: true }}
        onCreated={({ gl, scene }) => {
          scene.background = null
          gl.setClearColor(0x000000, 0)
        }}
      >
        <Suspense fallback={null}>
          <CharacterScene />
        </Suspense>
      </Canvas> */}
    </div>
  )
}
