'use client'

import { useEffect, useMemo, useRef } from 'react'
import { useFrame, useThree } from '@react-three/fiber'
import { useGLTF } from '@react-three/drei'
import * as THREE from 'three'

import { heroCharacterConfig } from './hero-character-config'

type GltfHeroCharacterProps = {
  url?: string
  lookStrength?: number
}

export function GltfHeroCharacter({
  url = heroCharacterConfig.gltfUrl,
  lookStrength = 0.35,
}: GltfHeroCharacterProps) {
  const root = useRef<THREE.Group>(null)
  const { scene, animations } = useGLTF(url)
  const pointer = useThree((state) => state.pointer)
  const mixerRef = useRef<THREE.AnimationMixer | null>(null)

  const cloned = useMemo(() => scene.clone(true), [scene])

  const parts = useMemo(
    () => ({
      torso: cloned.getObjectByName('Torso'),
      head: cloned.getObjectByName('Head'),
      leftArm: cloned.getObjectByName('LeftArm'),
      rightArm: cloned.getObjectByName('RightArm'),
    }),
    [cloned],
  )

  useEffect(() => {
    cloned.traverse((obj) => {
      const mesh = obj as THREE.Mesh
      if (mesh.isMesh && mesh.material) {
        const mats = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
        mats.forEach((m) => {
          const std = m as THREE.MeshStandardMaterial
          if (std.envMapIntensity !== undefined) std.envMapIntensity = 0.35
        })
      }
    })
  }, [cloned])

  useEffect(() => {
    const mixer = new THREE.AnimationMixer(cloned)
    mixerRef.current = mixer
    if (animations[0]) mixer.clipAction(animations[0]).play()
    return () => {
      mixer.stopAllAction()
      mixerRef.current = null
    }
  }, [cloned, animations])

  useFrame((state, delta) => {
    mixerRef.current?.update(delta)
    if (animations.length > 0) return

    const t = state.clock.elapsedTime
    if (root.current) {
      root.current.position.y = Math.sin(t * 1.25) * 0.025
      root.current.rotation.y = Math.sin(t * 0.45) * 0.05
    }
    if (parts.torso) parts.torso.scale.y = 1 + Math.sin(t * 1.25) * 0.01
    if (parts.leftArm) {
      parts.leftArm.rotation.x = Math.sin(t * 1.1) * 0.08 - 0.08
      parts.leftArm.rotation.z = 0.12
    }
    if (parts.rightArm) {
      parts.rightArm.rotation.x = Math.sin(t * 1.1 + Math.PI) * 0.08 - 0.06
      parts.rightArm.rotation.z = -0.12
    }
    if (parts.head) {
      parts.head.rotation.y = THREE.MathUtils.lerp(
        parts.head.rotation.y,
        pointer.x * lookStrength,
        0.08,
      )
      parts.head.rotation.x = THREE.MathUtils.lerp(
        parts.head.rotation.x,
        -pointer.y * lookStrength * 0.4,
        0.08,
      )
    }
  })

  return (
    <group ref={root}>
      <primitive object={cloned} />
    </group>
  )
}

if (heroCharacterConfig.mode === 'gltf') {
  useGLTF.preload(heroCharacterConfig.gltfUrl)
}
